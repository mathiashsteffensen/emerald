package evaluator

import (
	"emerald/object"
	"testing"
)

func TestMethodLiteral(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		method    string
		definedOn string
		isStatic  bool
	}{
		{
			"defining on main Object",
			`
			def log(msg)
				puts(msg)
			end`,
			"log",
			"Object",
			true,
		},
		{
			"overloading Integer class",
			`
			class Integer
				def times(x, y)
					x * y
				end
			end`,
			"times",
			"Integer",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			testSymbolObject(t, evaluated, tt.method)

			class, ok := object.NewEnvironment().Get(tt.definedOn)
			if !ok {
				t.Fatalf("class %s does not exist in the environment", tt.definedOn)
			}

			var target object.EmeraldValue

			if tt.isStatic {
				target = class
			} else {
				target = class.(*object.Class).New()
			}

			if !target.RespondsTo(tt.method, target) {
				t.Fatalf("%s does not respond to %s", tt.definedOn, tt.method)
			}
		})
	}
}
