package emerald_test

import (
	"emerald"
	"emerald/object"
	"fmt"
	"testing"
)

func TestEngineExceptionClassName(t *testing.T) {
	for _, parent := range []string{"Exception", "StandardError", "ArgumentError", "TypeError", "NameError", "NoMethodError", "RuntimeError", "LoadError"} {
		for _, tt := range []struct {
			name, source, className, message string
		}{
			{"owner", `raise ` + parent + `, "boom"`, parent, "boom"},
			{"subclass", fmt.Sprintf(`class Custom < %s; end; raise Custom, "boom"`, parent), "Custom", ""},
			{"descendant with singleton", fmt.Sprintf(`class Custom < %s; end
				class Child < Custom
					def initialize(message)
						class << self; def marker; 42; end; end
					end
				end
				raise Child, "boom"`, parent), "Child", ""},
		} {
			t.Run(parent+"/"+tt.name, func(t *testing.T) {
				_, err := emerald.New().Eval(tt.source)
				evalErr, ok := err.(emerald.EvalError)
				if !ok || evalErr.ClassName != tt.className || evalErr.Message != tt.message {
					t.Fatalf("error = %#v, want %s: %s", err, tt.className, tt.message)
				}
			})
		}
	}
}

func TestEngineEscapedProcKeepsYield(t *testing.T) {
	engine := emerald.New()

	_, err := engine.Eval(`
		def capture(&block)
			block
		end
		def outer
			capture { yield }
		end
		$saved = outer { 42 }
	`)
	if err != nil {
		t.Fatalf("first eval: %v", err)
	}

	value, err := engine.Eval(`$saved.call`)
	if err != nil {
		t.Fatalf("second eval: %v", err)
	}

	if !value.Is(object.INTEGER_VALUE) || value.Num != 42 {
		t.Fatalf("escaped Proc#call = %s, want 42", value.Inspect())
	}
}

func TestEngineEvalErrorBacktrace(t *testing.T) {
	_, err := emerald.New().EvalFile("main.rb", "def inner\n  raise \"boom\"\nend\n\ndef outer\n  inner\nend\n\nouter\n")
	evalErr, ok := err.(emerald.EvalError)
	if !ok {
		t.Fatalf("expected EvalError, got %T: %v", err, err)
	}

	wantBacktrace := []string{"main.rb:2:in 'inner'", "main.rb:6:in 'outer'", "main.rb:9:in '<main>'"}
	if fmt.Sprint(evalErr.Backtrace) != fmt.Sprint(wantBacktrace) {
		t.Fatalf("backtrace = %q, want %q", evalErr.Backtrace, wantBacktrace)
	}

	wantFull := "main.rb:2:in 'inner': boom (RuntimeError)\n\tfrom main.rb:6:in 'outer'\n\tfrom main.rb:9:in '<main>'"
	if evalErr.FullMessage() != wantFull {
		t.Fatalf("full message = %q, want %q", evalErr.FullMessage(), wantFull)
	}

	if evalErr.Error() != "RuntimeError: boom" {
		t.Fatalf("Error() = %q", evalErr.Error())
	}

	if (emerald.EvalError{ClassName: "Exception", Message: "no frames"}).FullMessage() != "no frames (Exception)" {
		t.Fatal("full message without a backtrace")
	}
}
