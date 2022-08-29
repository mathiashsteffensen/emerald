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
			input:    `(/a/ =~ "a").is_a?(MatchData)`,
			expected: true,
		},
		{
			input:    `/a/.match("a").is_a?(MatchData)`,
			expected: true,
		},
		{
			input:    `/a/.match("b")`,
			expected: nil,
		},
		{
			input: `
				/a/ =~ "a"
				$~.is_a?(MatchData)
			`,
			expected: true,
		},
		{
			input: `
				/a/ =~ "a"
				Regexp.last_match.is_a?(MatchData)
			`,
			expected: true,
		},
		{
			input: `
				/a/ =~ "a"
				$&
			`,
			expected: "a",
		},
		{
			input: `
				/a(b)(c)/ =~ "abc"
				$1
			`,
			expected: "b",
		},
		{
			input: `
				/a(b)(c)/ =~ "abc"
				$2
			`,
			expected: "c",
		},
	}

	runCoreTests(t, tests)
}
