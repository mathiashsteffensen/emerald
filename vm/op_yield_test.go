package vm

import "testing"

func TestOpYield(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				def run_with_num
					yield 5
				end

				run_with_num { |n| n + 5 }
			`,
			expected: 10,
		},
	}

	runVmTests(t, tests)
}
