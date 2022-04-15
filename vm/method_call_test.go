package vm

import "testing"

func TestMethodCall(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "built in instance method",
			input:    `"string".to_sym`,
			expected: ":string",
		},
		{
			name: "creating instance and instance method",
			input: `
			class Greeter
				def hello
					"hello"
				end
			end

			instance = Greeter.new

			instance.hello
			`,
			expected: "hello",
		},
	}

	runVmTests(t, tests)
}
