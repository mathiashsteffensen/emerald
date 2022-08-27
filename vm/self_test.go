package vm

import "testing"

func TestSelf(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "assignment call to self",
			input: `
				def value=(new)
					@value = new
				end

				self.value = 5
				value = 10

				@value
			`,
			expected: 5,
		},
	}

	runVmTests(t, tests)
}
