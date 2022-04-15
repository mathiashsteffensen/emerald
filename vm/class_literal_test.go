package vm

import "testing"

func TestClassLiteral(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "empty class",
			input:    "class MyClass; end",
			expected: nil,
		},
		{
			name: "class with method",
			input: `class MyClass
				def my_method; end
			end`,
			expected: ":my_method",
		},
		{
			name: "class with arbitrary last expression",
			input: `class MyClass
				def my_method; end
				def my_other_method; end

				69
			end`,
			expected: 69,
		},
	}

	runVmTests(t, tests)
}
