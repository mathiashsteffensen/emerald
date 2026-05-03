package core_test

import "testing"

func TestArray_compactBang(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "arr = [2, 3, nil, 5]; arr.compact!; arr",
			expected: []any{2, 3, 5},
		},
		{
			input:    "arr = [nil, 2, nil]; arr.compact!; arr",
			expected: []any{2},
		},
	}

	runCoreTests(t, tests)
}

func TestArray_push(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "arr = []; arr.push(2, 3, 4); arr",
			expected: []any{2, 3, 4},
		},
		{
			input:    "arr = []; arr << 2; arr",
			expected: []any{2},
		},
	}

	runCoreTests(t, tests)
}

func TestArray_pop(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "[2, 3, 4].pop",
			expected: 4,
		},
		{
			input:    "arr = [2, 3, 4]; arr.pop; arr",
			expected: []any{2, 3},
		},
		{
			input:    "[].pop",
			expected: nil,
		},
	}

	runCoreTests(t, tests)
}

func TestArray_to_s(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "[2, 3, \"boop\"].to_s",
			expected: "[2, 3, \"boop\"]",
		},
	}

	runCoreTests(t, tests)
}

func TestArray_first(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "[12, 24, 49].first",
			expected: 12,
		},
		{
			input:    "[36, 24, 49].first(2)",
			expected: []any{36, 24},
		},
		{
			input:    "[36, 24, 49].first(2, 3)",
			expected: "error:ArgumentError:wrong number of arguments (given 2, expected 1)",
		},
	}

	runCoreTests(t, tests)
}

func TestArray_find(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "[12, 24, 49].find { |i| i / 6 / 2 == 2 }",
			expected: 24,
		},
		{
			input:    "[12, 24, 49].find { |i| i / 6 / 2 == 3 }",
			expected: nil,
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

func TestArray_reduce(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "with just a block",
			input:    "[1, 3, 5].reduce { |sum, n| sum + n }",
			expected: 9,
		},
		{
			name:     "with an initial value argument & a block",
			input:    "[1, 3, 5].reduce 50 { |sum, n| sum + n }",
			expected: 59,
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

func TestArray_equals(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "[2, 4, 8] == [2, 4, 8]",
			expected: true,
		},
		{
			input:    "[2, 4, 8] == 2",
			expected: false,
		},
		{
			input:    "[2, 4, 8] == [2, 4]",
			expected: false,
		},
		{
			input:    "[2, 4, 8] == [2, 4, 9]",
			expected: false,
		},
	}

	runCoreTests(t, tests)
}
