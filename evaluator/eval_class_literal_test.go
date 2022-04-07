package evaluator

import (
	"emerald/object"
	"testing"
)

func TestEvalClassLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		class    string
		method   string
		isStatic bool
	}{
		{
			"extending a builtin class",
			`class Integer
				def do_math
					28 * 19
				end
			end`,
			"Integer",
			"do_math",
			false,
		},
		{
			"constructing a new class with an instance method",
			`class Logger
				def info(msg)
					puts("INFO | " + msg)
				end
			end`,
			"Logger",
			"info",
			false,
		},
		{
			"constructing a new class with a static method - eigenclass style",
			`class Logger
				class << self
					def info(msg)
						puts("INFO | " + msg)
					end
				end
			end`,
			"Logger",
			"info",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			testSymbolObject(t, evaluated, tt.method)

			class, ok := object.NewEnvironment().Get(tt.class)
			if !ok {
				t.Fatalf("class %s does not exist in the environment", tt.class)
			}

			var target object.EmeraldValue

			if tt.isStatic {
				target = class
			} else {
				target = class.(*object.Class).New()
			}

			if !target.RespondsTo(tt.method, target) {
				t.Fatalf("instance of %s does not respond to %s", tt.class, tt.method)
			}
		})
	}
}
