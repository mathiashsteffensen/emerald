package evaluator

import (
	"emerald/object"
	"testing"
)

func TestEvalClassLiteral(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		class  string
		method string
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
		},
		{
			"constructing a new class",
			`class Logger
				def info(msg)
					puts("INFO | " + msg)
				end
			end`,
			"Logger",
			"info",
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

			instance := class.(*object.Class).New()

			if !instance.RespondsTo(tt.method, instance) {
				t.Fatalf("instance of %s does not respond to %s", tt.class, tt.method)
			}
		})
	}
}
