package emerald_test

import (
	"context"
	"emerald"
	"emerald/object"
	"errors"
	"testing"
	"time"
)

func TestNewSandboxRequiresPositiveTimeout(t *testing.T) {
	if _, err := emerald.NewSandbox(emerald.SandboxOptions{}); err == nil {
		t.Fatal("expected an error for a missing timeout")
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
				Timeout: 20 * time.Millisecond,
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
			if elapsed := time.Since(start); elapsed > time.Second {
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

func newTestSandbox(t *testing.T) *emerald.Sandbox {
	t.Helper()

	sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout: 500 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	return sandbox
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
