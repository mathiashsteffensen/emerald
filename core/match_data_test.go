package core_test

import "testing"

func TestMatchData_index_accessor(t *testing.T) {
	tests := []coreTestCase{
		{
			input: `
				/a(b)/ =~ "ab"
				
				[$~[0], $~[1]]
			`,
			expected: []any{"ab", "b"},
		},
	}

	runCoreTests(t, tests)
}

func TestMatchData_to_s(t *testing.T) {
	tests := []coreTestCase{
		{
			input: `
				/a(b)/ =~ "ab"
				
				$~.to_s
			`,
			expected: "ab",
		},
	}

	runCoreTests(t, tests)
}

func TestMatchData_to_a(t *testing.T) {
	tests := []coreTestCase{
		{
			input: `
				m = /(.)(.)(\d+)(\d)/.match("THX1138.")
				m.to_a
			`,
			expected: []any{"HX1138", "H", "X", "113", "8"},
		},
	}

	runCoreTests(t, tests)
}

func TestMatchData_captures(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    `(/a/ =~ "a").captures`,
			expected: []any{},
		},
		{
			input:    `(/a(b)/ =~ "ab").captures`,
			expected: []any{"b"},
		},
	}

	runCoreTests(t, tests)
}

func TestMatchData_regexp(t *testing.T) {
	tests := []coreTestCase{
		{
			input: `
				re = /(.)(.)(\d+)(\d)/
				re.match("THX1138.").regexp == re
			`,
			expected: true,
		},
	}

	runCoreTests(t, tests)
}
