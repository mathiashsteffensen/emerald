package emerald_test

import (
	"emerald"
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
