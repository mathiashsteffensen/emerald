package core_test

import "testing"

func TestHash_index_setter(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "{}[:key] = :value",
			expected: ":value",
		},
		{
			input:    "hash = {}; hash[:key] = :value; hash[:key]",
			expected: ":value",
		},
	}

	runCoreTests(t, tests)
}

func TestHash_index_accessor(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "{ key: :value }[:key]",
			expected: ":value",
		},
		{
			input:    "{}[:key]",
			expected: nil,
		},
		{
			input:    "{}[]",
			expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1)",
		},
	}

	runCoreTests(t, tests)
}

func TestHash_map(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "when block takes a single argument",
			input:    "{key: 2}.map { |key| key.to_s }",
			expected: []any{"key"},
		},
		{
			name:     "when block takes multiple arguments",
			input:    "{key: 2}.map { |key, value| [key, value] }",
			expected: []any{[]any{":key", 2}},
		},
	}

	runCoreTests(t, tests)
}

func TestHash_find(t *testing.T) {
	cases := []coreTestCase{
		{
			input:    "{key: 2, \"key\" => 4}.find { |key| key == :key }",
			expected: []any{":key", 2},
		},
		{
			input:    "{key: 2, \"key\" => 4}.find { |key| key == \"key\" }",
			expected: []any{"key", 4},
		},
	}

	runCoreTests(t, cases)
}

func TestHash_equals(t *testing.T) {
	cases := []coreTestCase{
		{
			input:    "{} == {}",
			expected: true,
		},
		{
			input:    "{} == 2",
			expected: false,
		},
		{
			input:    "{\"key\" => 2} == { key: 2 }",
			expected: false,
		},
		{
			input:    "{ key: 2 } == { key: 2 }",
			expected: true,
		},
		{
			input:    "{ key: 2, other_key: 3 } == { key: 2, other_key: 4 }",
			expected: false,
		},
	}

	runCoreTests(t, cases)
}
