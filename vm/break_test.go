package vm

import "testing"

func TestBreak(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "in a while loop with a value",
			input: `
				i = 0
				while true
					i += 1
					break i * 10 if i == 3
				end
			`,
			expected: 30,
		},
		{
			name: "in a while loop without a value",
			input: `
				result = []
				items = [1, 2, 3]
				while item = items.pop
					break if item == 2
					result.push(item)
				end
				result
			`,
			expected: []any{3},
		},
		{
			name: "in a nested while loop",
			input: `
				result = []
				i = 0
				while i < 2
					i += 1
					while true
						break
					end
					result.push(i)
				end
				result
			`,
			expected: []any{1, 2},
		},
		{
			name:     "in a block given to a native method",
			input:    `[1, 2, 3].each { |i| break i * 100 if i == 2 }`,
			expected: 200,
		},
		{
			name:     "in a block without a value",
			input:    `[1, 2, 3].map { |i| break if i == 2; i }`,
			expected: nil,
		},
		{
			name: "in a block given to a bytecode method",
			input: `
				class Foo
					def each
						yield 1
						yield 2
						$reached = true
					end
				end
				[Foo.new.each { |i| break i + 40 if i == 2 }, $reached]
			`,
			expected: []any{42, nil},
		},
		{
			name: "in a nested block only exits the inner call",
			input: `
				[[1, 2], [3, 4]].map { |pair| pair.each { |i| break i } }
			`,
			expected: []any{1, 3},
		},
		{
			name: "inside a method keeps running the method",
			input: `
				def foo
					result = [1, 2, 3].each { |i| break :early if i == 2 }
					[result, :after]
				end
				foo
			`,
			expected: []any{":early", ":after"},
		},
		{
			name: "in a while loop inside a block",
			input: `
				[1].map do |i|
					while true
						break
					end
					i
				end
			`,
			expected: []any{1},
		},
		{
			name:     "at the top level",
			input:    `break`,
			expected: "error:Exception:Invalid break",
		},
		{
			name: "in a method body",
			input: `
				def foo
					break
				end
			`,
			expected: "error:Exception:Invalid break",
		},
	}

	runVmTests(t, tests)
}
