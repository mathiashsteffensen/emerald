package core_test

import "testing"

func TestHash_map(t *testing.T) {
	cases := []coreTestCase{
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

	runCoreTests(t, cases)
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
