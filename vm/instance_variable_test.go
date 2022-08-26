package vm

import "testing"

func TestInstanceVariable(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "setting an instance var on main object",
			input:    "@var = 5",
			expected: 5,
		},
		{
			name:     "getting an instance var from main object",
			input:    "@var = 5; @var + 10",
			expected: 15,
		},
		{
			name:     "setting an instance var on main object with boolean and",
			input:    "@var = true; @var &&= 5",
			expected: 5,
		},
		{
			name:     "failing with setting an instance var on main object with boolean and",
			input:    "@var = false; @var &&= 5",
			expected: false,
		},
		{
			name:     "setting an instance var on main object with boolean or",
			input:    "@var = false; @var ||= 5",
			expected: 5,
		},
		{
			name:     "failing with setting an instance var on main object with boolean and",
			input:    "@var = true; @var ||= 5",
			expected: true,
		},
		{
			name: "setting an instance var on custom class",
			input: `
				class MyClass
					def set
						@info = :info
					end
					def info
						@info
					end
				end
		
				instance = MyClass.new
				instance.set
				instance.info
			`,
			expected: ":info",
		},
		{
			name: "setting an instance var on custom class with boolean or",
			input: `
				class MyClass
					def num; 5; end

					class << self
						def instance; @instance ||= new; end
					end
				end

				MyClass.instance
				MyClass.instance.num
			`,
			expected: 5,
		},
	}

	runVmTests(t, tests)
}
