package vm

import "testing"

func TestNext(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "in a block with a value",
			input:    `[1, 2, 3].map { |i| next i * 10 if i == 2; i }`,
			expected: []any{1, 20, 3},
		},
		{
			name: "in a block without a value",
			input: `
				[1, 2].map do |i|
					next
					i
				end
			`,
			expected: []any{nil, nil},
		},
		{
			name: "in a while loop",
			input: `
				result = []
				i = 0
				while i < 5
					i += 1
					next if i == 2 || i == 4
					result.push(i)
				end
				result
			`,
			expected: []any{1, 3, 5},
		},
		{
			name: "in a while loop inside a case",
			input: `
				result = []
				items = [1, 2, 3]
				while item = items.pop
					case item
					when 2
						next
					end
					result.push(item)
				end
				result
			`,
			expected: []any{3, 1},
		},
		{
			name: "in a while loop inside a block",
			input: `
				[[1, 2], [3, 4]].map do |items|
					sum = 0
					while item = items.pop
						next if item == 2
						sum += item
					end
					sum
				end
			`,
			expected: []any{1, 7},
		},
		{
			name: "in a block inside a while loop",
			input: `
				result = []
				items = [1]
				while items.pop
					result = [1, 2].map { |i| next :skipped if i == 1; i }
				end
				result
			`,
			expected: []any{":skipped", 2},
		},
		{
			name: "in a for loop",
			input: `
				odd = []
				for i in [1, 2, 3]
					next if i == 2
					odd << i
				end
				odd
			`,
			expected: []any{1, 3},
		},
		{
			name: "in a for loop inside a while loop",
			input: `
				seen = []
				i = 0
				while i < 2
					i += 1
					for item in [1, 2]
						next if item == 1
						seen << [i, item]
					end
				end
				seen
			`,
			expected: []any{[]any{1, 2}, []any{2, 2}},
		},
		{
			name: "in a while loop inside a for loop",
			input: `
				seen = []
				for item in [1, 2]
					j = 0
					while j < 3
						j += 1
						next if j == 2
						seen << [item, j]
					end
				end
				seen
			`,
			expected: []any{[]any{1, 1}, []any{1, 3}, []any{2, 1}, []any{2, 3}},
		},
		{
			name:     "at the top level",
			input:    `next`,
			expected: "error:Exception:Invalid next",
		},
		{
			name: "in a method body",
			input: `
				def foo
					next
				end
			`,
			expected: "error:Exception:Invalid next",
		},
	}

	runVmTests(t, tests)
}
