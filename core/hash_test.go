package core_test

import "testing"

func TestHash_map(t *testing.T) {
	cases := []coreTestCase{
		{
			input:    "{key: 2}.map { |key| key.to_s }",
			expected: []any{"key"},
		},
	}

	runCoreTests(t, cases)
}
