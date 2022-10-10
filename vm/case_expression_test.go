package vm

import "testing"

func TestCaseExpression(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				case 2
				when Integer
					1
				else
					3
				end
			`,
			expected: 1,
		},
		{
			input: `
				case 2
				when String
					1
				else
					3
				end
			`,
			expected: 3,
		},
		{
			input: `
				case 11
				when String
					1
				when Integer
					2
				else
					3
				end
			`,
			expected: 2,
		},
	}

	runVmTests(t, tests)
}
