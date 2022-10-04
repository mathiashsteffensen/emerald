package vm

import "testing"

func TestOpBang(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    "!1",
			expected: false,
		},
		{
			input:    "!nil",
			expected: true,
		},
		{
			input: `
				class BooleanNegatable
					def !@
						5
					end
				end

				n = BooleanNegatable.new
				!n
			`,
			expected: 5,
		},
	}

	runVmTests(t, tests)
}
