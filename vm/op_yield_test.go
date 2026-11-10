package vm

import "testing"

func TestOpYield(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "repeated yield updates arguments",
			input: `def twice(n); yield(n) + yield(n + 1); end
				twice(3) { |n| n * 2 }`,
			expected: 14,
		},
		{
			name:     "yield leaves method arguments intact",
			input:    `def f(n); yield(n + 1); n; end; f(3) { |n| n * 2 }`,
			expected: 3,
		},
		{
			name:     "missing block parameters become nil",
			input:    `def f; yield; end; f { |a, b| [a, b] }`,
			expected: []any{nil, nil},
		},
		{
			name:     "extra yield arguments are ignored",
			input:    `def f; yield 1, 2; end; f { |a| a }`,
			expected: 1,
		},
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
