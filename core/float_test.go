package core_test

import "testing"

func TestFloat_to_s(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "3.1415.to_s",
			expected: "3.1415",
		},
	}

	runCoreTests(t, tests)
}

func TestFloat_infix_operators(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "3.1 <=> 4.2",
			expected: -1,
		},
		{
			input:    "5.1 <=> 4.2",
			expected: 1,
		},
		{
			input:    "4.2 <=> 4.2",
			expected: 0,
		},
		{
			input:    "3.1 <=> nil",
			expected: nil,
		},
		{
			input:    "3.1 + 4.2",
			expected: 7.3,
		},
		{
			input:    "3.1 + 4",
			expected: 7.1,
		},
		{
			input:    "4.2 - 3.1",
			expected: 1.1,
		},
		{
			input:    "4.2 - 3",
			expected: 1.2,
		},
		{
			input:    "2.2 * 1.1",
			expected: 2.42,
		},
		{
			input:    "2.2 * 2",
			expected: 4.4,
		},
		{
			input:    "2.2 / 1.1",
			expected: 2.0,
		},
	}

	runCoreTests(t, tests)
}
