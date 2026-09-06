package emerald_test

import (
	"context"
	"emerald"
	"emerald/object"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReadmeLanguageFeatures(t *testing.T) {
	path := "scripts/readme_features.rb"
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	value, err := emerald.New().EvalFile(path, string(source))
	if err != nil {
		t.Fatal(err)
	}
	if !value.Is(object.INTEGER_VALUE) || value.Num != 216 {
		t.Fatalf("feature script did not complete all 216 checks: %s", value.Inspect())
	}
}

func TestReadmeEngineStateAndInputs(t *testing.T) {
	engine := emerald.New()
	if _, err := engine.Eval(`$readme = 41
		class ReadmeState; def value; 42; end; end
		def readme_method; 43; end`); err != nil {
		t.Fatal(err)
	}
	value, err := engine.Eval(`[$readme, ReadmeState.new.value, readme_method] == [41, 42, 43]`)
	if err != nil || !value.Is(object.TRUE_VALUE) {
		t.Fatalf("state did not survive evaluation: %s, %v", value.Inspect(), err)
	}
	value, err = engine.EvalWithOptions(`[ARGV[0], $*[0], $LOAD_PATH[0]] == ["one", "one", "/app/lib"]`, emerald.EvalOptions{
		Args: []string{"one"}, LoadPath: []string{"/app/lib"},
	})
	if err != nil || !value.Is(object.TRUE_VALUE) {
		t.Fatalf("explicit inputs: %s, %v", value.Inspect(), err)
	}
	value, err = engine.Eval(`[ARGV.length, $*.length, $LOAD_PATH.length] == [0, 0, 0]`)
	if err != nil || !value.Is(object.TRUE_VALUE) {
		t.Fatalf("inputs leaked to next evaluation: %s, %v", value.Inspect(), err)
	}
}

func TestReadmeSandboxSourceName(t *testing.T) {
	sandbox := newTestSandbox(t)
	name := filepath.Join(t.TempDir(), "does-not-exist.rb")
	value, err := sandbox.EvalFile(name, "21 * 2")
	if err != nil || value.Inspect() != "42" {
		t.Fatalf("EvalFile must evaluate supplied source without reading a file: %s, %v", value.Inspect(), err)
	}
	if _, err := sandbox.EvalFile(name, "def"); err == nil || !strings.Contains(err.Error(), name) {
		t.Fatalf("source name missing from diagnostic: %v", err)
	}
}

func TestReadmeSandboxRejectedValues(t *testing.T) {
	for _, source := range []string{"Object", "Object.new", "hash = {}; [hash, hash]", "hash = {}; hash[:self] = hash; hash"} {
		t.Run(source, func(t *testing.T) {
			if _, err := newTestSandbox(t).Eval(source); !errors.Is(err, emerald.ErrSandboxUnsupportedResult) {
				t.Fatalf("expected unsupported result, got %v", err)
			}
		})
	}
}

func TestReadmeSandboxClassAndMethodIsolation(t *testing.T) {
	sandbox := newTestSandbox(t)
	if _, err := sandbox.Eval("class ReadmePrivate; end; def readme_private; 42; end; nil"); err != nil {
		t.Fatal(err)
	}
	for _, source := range []string{"ReadmePrivate", "readme_private"} {
		if _, err := sandbox.Eval(source); err == nil {
			t.Fatalf("state leaked between sandbox evaluations: %s", source)
		}
	}
}

func TestReadmeSandboxStartupDeadline(t *testing.T) {
	worker := writeTestWorker(t, "#!/bin/sh\nwhile :; do :; done\n")
	sandbox, err := emerald.NewSandbox(emerald.SandboxOptions{
		Timeout: 100 * time.Millisecond, MemoryLimitBytes: testSandboxMemoryLimitBytes, WorkerPath: worker,
	})
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	if _, err := sandbox.Eval("1"); !errors.Is(err, emerald.ErrSandboxTimeout) {
		t.Fatalf("startup did not match ErrSandboxTimeout: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("startup exceeded complete-evaluation deadline: %s", elapsed)
	}
}

func TestReadmeBuildFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte("#!/bin/sh\nexit 23\n"), 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("bash", "scripts/build", "emerald")
	cmd.Env = append(os.Environ(), "PATH="+dir+":"+os.Getenv("PATH"))
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("build reported success after Go failed:\n%s", output)
	}
}

func TestReadmeCommands(t *testing.T) {
	for _, name := range []string{"emerald", "iem"} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			binary := filepath.Join(t.TempDir(), name)
			build := exec.CommandContext(ctx, "go", "build", "-o", binary, "./cmd/"+name)
			if output, err := build.CombinedOutput(); err != nil {
				t.Fatalf("build: %v\n%s", err, output)
			}
			cmd := exec.CommandContext(ctx, binary)
			want := []string{"=> 42", "=> 8", "See you next time!"}
			if name == "emerald" {
				cmd.Args = append(cmd.Args, "scripts/readme_features.rb")
				want = []string{"PASS: 216 README feature checks"}
			} else {
				cmd.Stdin = strings.NewReader("21 * 2\n$readme = 7\n$readme + 1\nexit\n")
			}
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command: %v\n%s", err, output)
			}
			for _, marker := range want {
				if !strings.Contains(string(output), marker) {
					t.Fatalf("missing %q in output:\n%s", marker, output)
				}
			}
			if name == "emerald" {
				path := filepath.Join(t.TempDir(), "failure.rb")
				if err := os.WriteFile(path, []byte(`raise "readme failure"`), 0600); err != nil {
					t.Fatal(err)
				}
				output, err := exec.CommandContext(ctx, binary, path).CombinedOutput()
				if err == nil || !strings.Contains(string(output), "readme failure") {
					t.Fatalf("CLI must report failure with a nonzero exit: %v\n%s", err, output)
				}
			}
		})
	}
}
