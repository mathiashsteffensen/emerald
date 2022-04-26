package vm

import "testing"

func TestIndexAccessor(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "array index accessor",
			input:    "[0,1,2][1]",
			expected: 1,
		},
	}

	runVmTests(t, tests)
}
