package vm

import "testing"

func TestIfExpression(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				if true
					2
				elsif false
					3
				end
			`,
			expected: 2,
		},
		{
			input: `
				if false
					2
				elsif true
					3
				end
			`,
			expected: 3,
		},
		{
			input: `
				if false
					2
				elsif false
					3
				else
					4
				end
			`,
			expected: 4,
		},
	}

	runVmTests(t, tests)
}
