package vm

import "testing"

func TestIfModifier(t *testing.T) {
	setup := `
		i = 1
		
		def call
			i = i + 1
		end
	`

	tests := []vmTestCase{
		{
			input: `
				call if true && false

				i
			`,
			expected: 1,
		},
		{
			input: `
				call if 1 + 1 == 2

				i
			`,
			expected: 2,
		},
		{
			input: `
				def call?(should_we)
					2 * 5 == 10 if should_we
				end

				call unless call? false
				call if call? 2 * 9 + 17 == 35
				call if call? false
				call if call? true ||
				
									false
				i
			`,
			expected: 4,
		},
	}

	runVmTests(t, tests, setup)
}
