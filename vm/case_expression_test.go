package vm

import "testing"

func TestCaseExpression(t *testing.T) {
	tests := []vmTestCase{
		{name: "regexp match", input: `case "abc"; when /b/; true; else; false; end`, expected: true},
		{name: "regexp miss", input: `/x/ === "abc"`, expected: false},
		{name: "regexp non-string", input: `/x/ === nil`, expected: false},
		{name: "semicolon before else body", input: "case 2; when 1; 3; else; 4; end", expected: 4},
		{name: "empty else body", input: "case 2; when 1; 3; else; end", expected: nil},
		{
			input: `
				case 2
				when Integer
					1
				else
					3
				end
			`,
			expected: 1,
		},
		{
			input: `
				case 2
				when String
					1
				else
					3
				end
			`,
			expected: 3,
		},
		{
			input: `
				case 11
				when String
					1
				when Integer
					2
				else
					3
				end
			`,
			expected: 2,
		},
	}

	runVmTests(t, tests)
}
