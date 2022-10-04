package vm

import "testing"

func TestOpConstant(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    "CONST = 5",
			expected: 5,
		},
		{
			input:    "CONST = 5; CONST",
			expected: 5,
		},
		{
			input: `
				class C
					CONST = 5
					class << self
						def get_const; CONST; end
					end
				end

				C.get_const
			`,
			expected: 5,
		},
	}

	runVmTests(t, tests)
}
