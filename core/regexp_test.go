package core_test

import "testing"

func TestRegexp_inspect(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "/abc/.inspect",
			expected: "/abc/",
		},
	}

	runCoreTests(t, tests)
}

func TestRegexp_match(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `/a/ =~ "a"`,
			expected: true,
		},
		{
			input:    `/a/ =~ "b"`,
			expected: false,
		},
		{
			input:    `/a/.match("a")`,
			expected: true,
		},
		{
			input:    `/a/.match("b")`,
			expected: false,
		},
	}

	runCoreTests(t, tests)
}
