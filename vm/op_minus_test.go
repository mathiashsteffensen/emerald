package vm

import "testing"

func TestOpMinus(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    "-1",
			expected: -1,
		},
		{
			input:    "-1.1",
			expected: -1.1,
		},
		{
			input: `
				class Negatable
					def -@
						5
					end
				end

				n = Negatable.new
				-n
			`,
			expected: 5,
		},
	}

	runVmTests(t, tests)
}
