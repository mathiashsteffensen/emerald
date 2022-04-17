package vm

import "testing"

func TestInstanceVariable(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "setting an instance var on main object",
			input:    "@var = 5",
			expected: 5,
		},

		{
			name:     "getting an instance var from main object",
			input:    "@var = 5; @var + 10",
			expected: 15,
		},
	}

	runVmTests(t, tests)
}
