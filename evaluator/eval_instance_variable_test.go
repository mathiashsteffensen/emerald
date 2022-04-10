package evaluator

import (
	"emerald/object"
	"testing"
)

func TestEvalInstanceVariable(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		definedOn string
		varName   string
		returns   any
	}{
		{
			"defining on main object",
			"@var = 5",
			"Object",
			"@var",
			5,
		},
		{
			"defining on an instance of a custom class",
			`
			class MyClass
				def set(val)
					@var = val
				end
			end

			instance = MyClass.new
			instance.set(15)`,
			"instance",
			"@var",
			15,
		},
		{
			"defining on and getting from an instance of a custom class",
			`
			class MyClass
				def set(val)
					@var = val
				end

				def get
					@var
				end
			end

			instance = MyClass.new
			instance.set(55)
			instance.get`,
			"instance",
			"@var",
			55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := object.NewEnvironment()

			evaluated := testEval(tt.input, env)

			testObjectValue(t, evaluated, tt.returns)

			target, ok := env.Get(tt.definedOn)
			if !ok {
				t.Fatalf("class %s does not exist in the environment", tt.definedOn)
			}

			_, isStatic := target.(*object.Class)

			val := target.InstanceVariableGet(isStatic, tt.varName, target, target)

			testObjectValue(t, val, tt.returns)
		})
	}
}
