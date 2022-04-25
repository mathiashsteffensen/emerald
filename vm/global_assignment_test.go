package vm

import "testing"

func TestGlobalAssignment(t *testing.T) {
	tests := []vmTestCase{
		{"", "one = 1", 1},
		{"", "one = 1; one", 1},
		{"", "one = 1; two = 2; one + two", 3},
		{"", "one = 1; two = one + one; one + two", 3},
		{"", "one = true; one && one = 15", 15},
		{"", "one = false; one || one = 15", 15},
		{"", "one = false; one && one = 15", false},
		{"", "one = true; one || one = 15", true},
		{"", "one = true; one &&= 15", 15},
		{"", "one = false; one ||= 15", 15},
		{"", "one = false; one &&= 15", false},
		{"", "one = true; one ||= 15", true},
		{
			input: `
				def method
					$var = "hello"
				end

				method

				$var
			`,
			expected: "hello",
		},
	}
	runVmTests(t, tests)
}
