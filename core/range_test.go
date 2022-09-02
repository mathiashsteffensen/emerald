package core_test

import "testing"

func TestRange_new(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "Range.new(1, 4)",
			expected: "instance:Range",
		},
	}

	runCoreTests(t, tests)
}

func TestRange_enumerable(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "Range#map",
			input:    "Range.new(1, 4).map { |n| n*2 }",
			expected: []any{2, 4, 6, 8},
		},
		{
			name:     "Range#reduce",
			input:    "Range.new(0, 8).reduce([1,0]) { |acc, w| [acc[1], acc[0]+acc[1]] }[0]",
			expected: 21,
		},
	}

	runCoreTests(t, tests)
}
