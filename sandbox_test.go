package emerald_test

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"emerald"
	"emerald/core"
	"emerald/object"
)

const testSandboxMemoryLimitBytes int64 = 256 * 1024 * 1024

var testSandboxWorkerPath string

func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "emerald-sandbox-worker-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create sandbox worker directory: %v\n", err)
		os.Exit(1)
	}

	testSandboxWorkerPath = filepath.Join(tempDir, "emerald-sandbox-worker")
	build := exec.Command("go", "build", "-o", testSandboxWorkerPath, "./cmd/emerald-sandbox-worker")
	if output, err := build.CombinedOutput(); err != nil {
		_ = os.RemoveAll(tempDir)
		fmt.Fprintf(os.Stderr, "build sandbox worker: %v\n%s", err, output)
		os.Exit(1)
	}

	exitCode := m.Run()
	if err := os.RemoveAll(tempDir); err != nil {
		fmt.Fprintf(os.Stderr, "remove sandbox worker directory: %v\n", err)
		if exitCode == 0 {
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func TestNewSandboxRequiresPositiveTimeout(t *testing.T) {
	if _, err := emerald.NewSandbox(emerald.SandboxOptions{
		MemoryLimitBytes: testSandboxMemoryLimitBytes,
		WorkerPath:       testSandboxWorkerPath,
	}); err == nil {
		t.Fatal("expected an error for a missing timeout")
	}
}

func TestNewSandboxRequiresMemoryLimitAndWorker(t *testing.T) {
	if _, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout:    time.Second,
		WorkerPath: testSandboxWorkerPath,
	}); err == nil {
		t.Fatal("expected an error for a missing memory limit")
	}
	if _, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout:          time.Second,
		MemoryLimitBytes: testSandboxMemoryLimitBytes,
	}); err == nil {
		t.Fatal("expected an error for a missing worker path")
	}
}

func TestSandboxEvaluatesPureCode(t *testing.T) {
	sandbox := newTestSandbox(t)

	value, err := sandbox.Eval("21 * 2")
	if err != nil {
		t.Fatalf("Eval returned an error: %s", err)
	}
	if !value.Is(object.INTEGER_VALUE) || int64(value.Num) != 42 {
		t.Fatalf("expected 42, got %s", value.Inspect())
	}
}

func TestSandboxRoundTripsCompositeResult(t *testing.T) {
	value, err := newTestSandbox(t).Eval(`[nil, true, false, 42, 1.5, "hello", :symbol, {"key" => [1, 2]}]`)
	if err != nil {
		t.Fatalf("Eval returned an error: %s", err)
	}

	array, ok := value.Heap.(*core.ArrayInstance)
	if !ok {
		t.Fatalf("expected an Array result, got %T (%s)", value.Heap, value.Inspect())
	}
	if len(array.Value) != 8 {
		t.Fatalf("expected 8 array values, got %d", len(array.Value))
	}

	for i, want := range []string{"nil", "true", "false", "42", "1.5", "hello", ":symbol"} {
		if got := array.Value[i].Inspect(); got != want {
			t.Fatalf("array value %d: expected %q, got %q", i, want, got)
		}
	}

	hash, ok := array.Value[7].Heap.(*core.HashInstance)
	if !ok {
		t.Fatalf("expected a Hash result, got %T", array.Value[7].Heap)
	}
	if hash.Values.Len() != 1 {
		t.Fatalf("expected one hash entry, got %d", hash.Values.Len())
	}
	pair := hash.Values.Front()
	if pair == nil {
		t.Fatal("expected a hash entry")
	}
	if got := hash.Keys[pair.Key].Inspect(); got != "key" {
		t.Fatalf("expected hash key %q, got %q", "key", got)
	}
	nested, ok := pair.Value.Heap.(*core.ArrayInstance)
	if !ok {
		t.Fatalf("expected a nested Array result, got %T", pair.Value.Heap)
	}
	if len(nested.Value) != 2 || nested.Value[0].Inspect() != "1" || nested.Value[1].Inspect() != "2" {
		t.Fatalf("expected nested [1, 2], got %#v", nested.Value)
	}
}

func TestSandboxRejectsUnsupportedRegexpResult(t *testing.T) {
	_, err := newTestSandbox(t).Eval("/emerald/")
	if !errors.Is(err, emerald.ErrSandboxUnsupportedResult) {
		t.Fatalf("expected unsupported result error, got %v", err)
	}
}

func TestSandboxRejectsAliasedAndCyclicResults(t *testing.T) {
	for name, source := range map[string]string{
		"alias": `
			a = [1]
			[a, a]
		`,
		"string alias": `
			s = "value"
			[s, s]
		`,
		"cycle": `
			a = []
			a << a
			a
		`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := newTestSandbox(t).Eval(source)
			if !errors.Is(err, emerald.ErrSandboxUnsupportedResult) {
				t.Fatalf("expected unsupported result error, got %v", err)
			}
		})
	}
}

func TestSandboxContainsWorkerFailure(t *testing.T) {
	failingWorker := writeTestWorker(t, "#!/bin/sh\nkill -TERM $$\n")
	sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout:          time.Second,
		MemoryLimitBytes: testSandboxMemoryLimitBytes,
		WorkerPath:       failingWorker,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := sandbox.Eval("1 + 1"); !errors.Is(err, emerald.ErrSandboxWorkerFailed) {
		t.Fatalf("expected worker failure, got %v", err)
	}

	value, err := newTestSandbox(t).Eval("21 * 2")
	if err != nil {
		t.Fatalf("subsequent Eval returned an error: %s", err)
	}
	if int64(value.Num) != 42 {
		t.Fatalf("expected 42, got %s", value.Inspect())
	}
}

func TestSandboxRejectsMalformedWorkerReply(t *testing.T) {
	malformedWorker := writeTestWorker(t, "#!/bin/sh\nprintf '\\000\\000\\000\\000'\nwhile :; do :; done\n")
	sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout:          5 * time.Second,
		MemoryLimitBytes: testSandboxMemoryLimitBytes,
		WorkerPath:       malformedWorker,
	})
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	_, err = sandbox.Eval("1")
	if !errors.Is(err, emerald.ErrSandboxWorkerFailed) {
		t.Fatalf("expected worker failure, got %v", err)
	}
	if elapsed := time.Since(start); elapsed >= 2*time.Second {
		t.Fatalf("malformed worker was not stopped promptly: %s", elapsed)
	}
}

func TestSandboxRejectsMismatchedWorkerVersion(t *testing.T) {
	mismatchedWorker := writeTestWorker(t, framedWorkerScript(`{"version":2,"value":{"type":"nil"}}`))
	sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout:          time.Second,
		MemoryLimitBytes: testSandboxMemoryLimitBytes,
		WorkerPath:       mismatchedWorker,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = sandbox.Eval("1")
	if !errors.Is(err, emerald.ErrSandboxWorkerFailed) {
		t.Fatalf("expected worker failure, got %v", err)
	}
	if !strings.Contains(err.Error(), "protocol version mismatch") {
		t.Fatalf("expected version mismatch detail, got %v", err)
	}
}

func TestSandboxRejectsOversizedWorkerReply(t *testing.T) {
	// 0x01000001 is one byte larger than the 16 MiB protocol ceiling.
	oversizedWorker := writeTestWorker(t, "#!/bin/sh\nprintf '\\001\\000\\000\\001'\n")
	sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout:          time.Second,
		MemoryLimitBytes: testSandboxMemoryLimitBytes,
		WorkerPath:       oversizedWorker,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = sandbox.Eval("1")
	if !errors.Is(err, emerald.ErrSandboxWorkerFailed) {
		t.Fatalf("expected worker failure, got %v", err)
	}
}

func TestSandboxCompilesTrailingCommentAtEOF(t *testing.T) {
	value, err := newTestSandbox(t).Eval("1 # trailing comment")
	if err != nil {
		t.Fatalf("Eval returned an error: %s", err)
	}
	if !value.Is(object.INTEGER_VALUE) || int64(value.Num) != 1 {
		t.Fatalf("expected 1, got %s", value.Inspect())
	}
}

func TestSandboxDoesNotDefineHostCapabilities(t *testing.T) {
	forbiddenConstants := []string{
		"IO",
		"File",
		"Dir",
		"Time",
		"TCPServer",
		"TCPSocket",
	}

	for _, constant := range forbiddenConstants {
		t.Run(constant, func(t *testing.T) {
			_, err := newTestSandbox(t).Eval(constant)
			assertEvalErrorClass(t, err, "NameError")
		})
	}

	forbiddenMethods := []string{
		`puts("hello")`,
		`print("hello")`,
		`require_relative("other")`,
		`sleep(1)`,
	}

	for _, script := range forbiddenMethods {
		t.Run(script, func(t *testing.T) {
			_, err := newTestSandbox(t).Eval(script)
			assertEvalErrorClass(t, err, "NoMethodError")
		})
	}
}

func TestSandboxDoesNotInheritHostInputs(t *testing.T) {
	sandbox := newTestSandbox(t)

	for _, script := range []string{"ARGV.size", "$*.size", "$LOAD_PATH.size"} {
		value, err := sandbox.Eval(script)
		if err != nil {
			t.Fatalf("%s returned an error: %s", script, err)
		}
		if !value.Is(object.INTEGER_VALUE) || int64(value.Num) != 0 {
			t.Fatalf("%s: expected 0, got %s", script, value.Inspect())
		}
	}
}

func TestEngineUsesExplicitHostInputs(t *testing.T) {
	engine := emerald.New()

	value, err := engine.Eval("ARGV.size")
	if err != nil {
		t.Fatalf("default Eval returned an error: %s", err)
	}
	if int64(value.Num) != 0 {
		t.Fatalf("expected default ARGV to be empty, got %s", value.Inspect())
	}

	value, err = engine.EvalWithOptions("ARGV.size + $LOAD_PATH.size", emerald.EvalOptions{
		Args:     []string{"one", "two"},
		LoadPath: []string{"/stdlib"},
	})
	if err != nil {
		t.Fatalf("EvalWithOptions returned an error: %s", err)
	}
	if int64(value.Num) != 3 {
		t.Fatalf("expected explicit inputs to have size 3, got %s", value.Inspect())
	}
}

func TestSandboxTimesOutBytecodeAndNativeLoops(t *testing.T) {
	tests := map[string]string{
		"bytecode": "while true\nend",
		"native":   "9223372036854775807.times { |i| i }",
		"nested": `
			def spin()
				while true
				end
			end
			spin
		`,
	}

	for name, script := range tests {
		t.Run(name, func(t *testing.T) {
			sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
				Timeout:          time.Second,
				MemoryLimitBytes: testSandboxMemoryLimitBytes,
				WorkerPath:       testSandboxWorkerPath,
			})
			if err != nil {
				t.Fatal(err)
			}

			start := time.Now()
			_, err = sandbox.Eval(script)
			if !errors.Is(err, emerald.ErrSandboxTimeout) {
				t.Fatalf("expected sandbox timeout, got %v", err)
			}
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected context deadline identity, got %v", err)
			}
			var timeoutErr *emerald.SandboxTimeoutError
			if !errors.As(err, &timeoutErr) || timeoutErr.Phase != "execute" {
				t.Fatalf("expected execute timeout details, got %#v", timeoutErr)
			}
			if elapsed := time.Since(start); elapsed > 3*time.Second {
				t.Fatalf("timeout took too long: %s", elapsed)
			}
		})
	}
}

func TestSandboxHonorsCallerCancellation(t *testing.T) {
	sandbox := newTestSandbox(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := sandbox.EvalContext(ctx, "1 + 2")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}

func TestSandboxPreservesCallerDeadline(t *testing.T) {
	sandbox := newTestSandbox(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, err := sandbox.EvalContext(ctx, "while true\nend")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected caller deadline, got %v", err)
	}
	if errors.Is(err, emerald.ErrSandboxTimeout) {
		t.Fatalf("caller deadline was misclassified as sandbox timeout: %v", err)
	}
}

func TestSandboxEvaluationsDoNotShareState(t *testing.T) {
	sandbox := newTestSandbox(t)

	if _, err := sandbox.Eval("$sandbox_state = 42"); err != nil {
		t.Fatalf("first Eval returned an error: %s", err)
	}

	value, err := sandbox.Eval("$sandbox_state")
	if err != nil {
		t.Fatalf("second Eval returned an error: %s", err)
	}
	if !value.IsNil() {
		t.Fatalf("expected isolated state, got %s", value.Inspect())
	}
}

func TestSandboxHandlesUserCreatedNamespaceAndIncludeCycles(t *testing.T) {
	tests := map[string]string{
		"namespace": `
			module A
				A = self
			end
			A.name
		`,
		"include": `
			module A
				include A
			end
			A.methods.size
		`,
	}

	for name, script := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := newTestSandbox(t).Eval(script); err != nil {
				t.Fatalf("Eval returned an error: %s", err)
			}
		})
	}
}

func TestSandboxSupportsConcurrentFreshEvaluations(t *testing.T) {
	sandbox := newTestSandbox(t)

	const evaluations = 8
	errs := make(chan error, evaluations)
	var group sync.WaitGroup
	for i := 0; i < evaluations; i++ {
		group.Add(1)
		go func(input int) {
			defer group.Done()

			value, err := sandbox.Eval(fmt.Sprintf("%d * %d", input, input))
			if err != nil {
				errs <- err
				return
			}
			if !value.Is(object.INTEGER_VALUE) || int64(value.Num) != int64(input*input) {
				errs <- fmt.Errorf("%d * %d = %s", input, input, value.Inspect())
			}
		}(i)
	}
	group.Wait()
	close(errs)

	for err := range errs {
		t.Fatal(err)
	}
}

func TestSandboxLinuxHardMemoryLimitTerminatesWorker(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("RLIMIT_AS is Linux-only")
	}

	sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout:          10 * time.Second,
		MemoryLimitBytes: testSandboxMemoryLimitBytes,
		WorkerPath:       testSandboxWorkerPath,
	})
	if err != nil {
		t.Fatal(err)
	}

	// A 1 MiB chunk reaches the 256 MiB address-space ceiling in a handful of
	// native string-builder growth steps, without placing that allocation in the
	// test process. The result is deliberately small so the transfer limit
	// cannot cause a false positive.
	source := fmt.Sprintf("value = %q * 384\n1", strings.Repeat("x", 1<<20))
	_, err = sandbox.Eval(source)
	if !errors.Is(err, emerald.ErrSandboxWorkerFailed) {
		t.Fatalf("expected hard-limit worker failure, got %v", err)
	}
}

func newTestSandbox(t *testing.T) *emerald.Sandbox {
	t.Helper()

	sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout:          5 * time.Second,
		MemoryLimitBytes: testSandboxMemoryLimitBytes,
		WorkerPath:       testSandboxWorkerPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	return sandbox
}

func writeTestWorker(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "worker")
	if err := os.WriteFile(path, []byte(content), 0o700); err != nil {
		t.Fatalf("write test worker: %v", err)
	}
	return path
}

func framedWorkerScript(payload string) string {
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))

	var escaped strings.Builder
	for _, b := range header {
		fmt.Fprintf(&escaped, "\\%03o", b)
	}
	return "#!/bin/sh\nprintf '" + escaped.String() + payload + "'\n"
}

func assertEvalErrorClass(t *testing.T, err error, className string) {
	t.Helper()

	var evalErr emerald.EvalError
	if !errors.As(err, &evalErr) {
		t.Fatalf("expected EvalError, got %T: %v", err, err)
	}
	if evalErr.ClassName != className {
		t.Fatalf("expected %s, got %s: %v", className, evalErr.ClassName, err)
	}
}
