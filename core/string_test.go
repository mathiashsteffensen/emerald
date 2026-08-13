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

func TestString_add(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"hello" + " " + "world"`,
			expected: "hello world",
		},
		{
			input:    `"wdaw" + 2`,
			expected: "error:TypeError:no implicit conversion of Integer into String",
		},
	}

	runCoreTests(t, tests)
}

func TestString_multiply(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"w" * 5`,
			expected: "wwwww",
		},
		{
			input:    `"w" * ""`,
			expected: "error:TypeError:no implicit conversion of String into Integer",
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

func TestString_size(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"Hello".size`,
			expected: 5,
		},
	}

	runCoreTests(t, tests)
}

func TestString_split(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `"h e llo o".split`,
			expected: []any{"h", "e", "llo", "o"},
		},
		{
			input:    `"h e llo o".split("o")`,
			expected: []any{"h e ll", " ", ""},
		},
	}

	runCoreTests(t, tests)
}

func TestString_start_with(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "matches an exact prefix",
			input:    `"hello".start_with?("hello")`,
			expected: true,
		},
		{
			name:     "matches a shorter prefix",
			input:    `"hello".start_with?("he")`,
			expected: true,
		},
		{
			name:     "does not match text occurring later",
			input:    `"hello".start_with?("ell")`,
			expected: false,
		},
		{
			name:     "does not match a prefix longer than the receiver",
			input:    `"hi".start_with?("hello")`,
			expected: false,
		},
		{
			name:     "matches an empty prefix",
			input:    `"hello".start_with?("")`,
			expected: true,
		},
		{
			name:     "does not match a nonempty prefix on an empty receiver",
			input:    `"".start_with?("hello")`,
			expected: false,
		},
		{
			name:     "returns false without prefixes",
			input:    `"hello".start_with?`,
			expected: false,
		},
		{
			name:     "matches any prefix argument",
			input:    `"hello".start_with?("no", "he")`,
			expected: true,
		},
		{
			name:     "returns false when no prefixes match",
			input:    `"hello".start_with?("no", "x")`,
			expected: false,
		},
		{
			name:     "short circuits after a matching prefix",
			input:    `"hello".start_with?("he", 1)`,
			expected: true,
		},
		{
			name:     "requires string prefixes until a match is found",
			input:    `"hello".start_with?("no", 1)`,
			expected: "error:TypeError:no implicit conversion of Integer into String",
		},
	}

	runCoreTests(t, tests)
}
