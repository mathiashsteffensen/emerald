package core_test

import "testing"

func TestArray_find(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "[12, 24, 49].find { |i| i / 6 / 2 == 2 }",
			expected: 24,
		},
	}

	runCoreTests(t, tests)
}

func TestArray_find_index(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "[12, 24, 49].find_index { |i| i / 6 / 2 == 2 }",
			expected: 1,
		},
	}

	runCoreTests(t, tests)
}

func TestArray_map(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "increment",
			input:    "[12, 24, 49].map { |i| i + 1 }",
			expected: []any{13, 25, 50},
		},
	}

	runCoreTests(t, tests)
}

func TestArray_sum(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "with no block or argument",
			input:    "[12, 24, 49].sum",
			expected: 85,
		},
		{
			name:     "with no block, but with argument",
			input:    "[12, 24, 49].sum(15)",
			expected: 100,
		},
		{
			name:     "with block, but no argument",
			input:    "[12, 24, 49].sum { |i| i * 10 }",
			expected: 850,
		},
		{
			name:     "with block and argument",
			input:    "[12, 24, 49].sum(150) { |i| i * 10 }",
			expected: 1000,
		},
	}

	runCoreTests(t, tests)
}

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
