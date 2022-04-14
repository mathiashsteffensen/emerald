package vm

import (
	"testing"
)

func TestMethodCall(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "five_plus_ten",
			input: `
			def five_plus_ten
				5 + 10
			end
			five_plus_ten
			`,
			expected: 15,
		},
		{
			name: "one + two",
			input: `
			def one; 1; end
			def two; 2; end
			one + two
			`,
			expected: 3,
		},
		{
			name: "1+1+1",
			input: `
			def a; 1; end
			def b; a + 1; end
			def c; b + 1; end
			c
			`,
			expected: 3,
		},
		{
			name: "early_exit",
			input: `
			def early_exit
				return 99
				100
			end
			early_exit
			`,
			expected: 99,
		},
		{
			name: "early_exit_two",
			input: `
			def early_exit
				return 99
				return 100
			end
			early_exit
			`,
			expected: 99,
		},
		{
			name: "no_return",
			input: `
			def no_return; end
			no_return
			`,
			expected: nil,
		},
		{
			name: "one_and_two",
			input: `
			def one_and_two
				one = 1
				two = 2
				one + two
			end
			one_and_two
			`,
			expected: 3,
		},
		{
			name: "one_and_two & three_and_four",
			input: `
			def one_and_two
				one = 1
				two = 2
				one + two
			end
			def three_and_four
				three = 3
				four = 4
				three + four
			end
			one_and_two + three_and_four
			`,
			expected: 10,
		},
		{
			input: `
			def first_foobar
				foobar = 50
				foobar
			end
			def second_foobar
				foobar = 100
				foobar
			end
			first_foobar + second_foobar
			`,
			expected: 150,
		},
		{
			input: `
			globalSeed = 50

			def minus_one
				num = 1
				globalSeed - num
			end

			def minus_two
				num = 2
				globalSeed - num
			end
			minus_one + minus_two
			`,
			expected: 97,
		},
		{
			input: `
			def identity(a)
				a
			end
			identity(4)
			`,
			expected: 4,
		},
		{
			input: `
			def sum(a, b)
				a + b
			end
			sum(1, 2)`,
			expected: 3,
		},
		{
			input: `
			def sum(a, b)
				c = a + b
				c
			end
			sum(1, 2)`,
			expected: 3,
		},
		{
			input: `
			def sum(a, b)
				c = a + b
				c
			end
			sum(1, 2) + sum(3, 4)`,
			expected: 10,
		},
		{
			input: `
			def sum(a, b)
				c = a + b
				c
			end

			def outer
				sum(1, 2) + sum(3, 4)
			end

			outer`,
			expected: 10,
		},
	}

	runVmTests(t, tests)
}
