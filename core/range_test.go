package core_test

import "testing"

func TestRange_new(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "rt.Range.new(1, 4)",
			expected: "instance:rt.Range",
		},
	}

	runCoreTests(t, tests)
}

func TestRange_enumerable(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "rt.Range#map",
			input:    "rt.Range.new(1, 4).map { |n| n*2 }",
			expected: []any{2, 4, 6, 8},
		},
		{
			name:     "rt.Range#reduce",
			input:    "rt.Range.new(0, 8).reduce([1,0]) { |acc, w| [acc[1], acc[0]+acc[1]] }[0]",
			expected: 21,
		},
	}

	runCoreTests(t, tests)
}
