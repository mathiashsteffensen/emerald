package core_test

import "testing"

func TestArray_index(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "element exists",
			input:    "[12, 24, 49][1]",
			expected: 24,
		},
		{
			name:     "element doesn't exist",
			input:    "[12, 24, 49][16]",
			expected: nil,
		},
		{
			name:     "wrong type argument",
			input:    `[12, 24, 49]["a string"]`,
			expected: "error:TypeError:no implicit conversion of String into Integer",
		},
	}

	runCoreTests(t, tests)
}
