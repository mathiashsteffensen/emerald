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
	}{
		{
			"defining on main Object",
			`
			def log(msg)
				puts(msg)
			end`,
			"log",
			"Object",
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

			if !class.RespondsTo(tt.method, class) {
				t.Fatalf("class %s does not respond to %s", tt.definedOn, tt.method)
			}
		})
	}
}
