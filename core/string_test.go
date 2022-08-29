package core_test

import "testing"

func TestString_to_sym(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"hello".to_sym`,
			expected: ":hello",
		},
	}

	runCoreTests(t, tests)
}

func TestString_match(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `("a" =~ /a/).is_a?(MatchData)`,
			expected: true,
		},
		{
			input:    `"a".match(/a/).is_a?(MatchData)`,
			expected: true,
		},
	}

	runCoreTests(t, tests)
}

func TestString_upcase(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"hello".upcase`,
			expected: "HELLO",
		},
	}

	runCoreTests(t, tests)
}
