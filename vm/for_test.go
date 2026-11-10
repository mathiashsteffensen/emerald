package vm

import "testing"

func TestFor(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "on multiple lines",
			input: `
				sum = 0
				for item in [1, 2, 3]
					sum += item
				end
				sum
			`,
			expected: 6,
		},
		{
			name:     "on one line",
			input:    "sum = 0; for item in [1, 2]; sum += item; end; sum",
			expected: 3,
		},
		{
			name:     "on one line with do",
			input:    "for item in [1, 2] do last = item * 10 end; last",
			expected: 20,
		},
		{
			name:     "on one line with do and an identifier ending the body",
			input:    "s = 0; for x in [1, 2] do s += x end; s",
			expected: 3,
		},
		{
			name:     "over a command call and a newline",
			input:    "def pair(a, b); [a, b]; end; seen = []\nfor x in pair 1, 2\nseen << x\nend\nseen",
			expected: []any{1, 2},
		},
		{
			name:     "over a command call and ';'",
			input:    "def pair(a, b); [a, b]; end; seen = []; for x in pair 3, 4; seen << x; end; seen",
			expected: []any{3, 4},
		},
		{
			name:     "over a collection with a do block in call parentheses",
			input:    "def wrap(x); x; end; for a in wrap([3].map do |v| v * 2 end) do last = a end; last",
			expected: 6,
		},
		{
			name:     "over a collection with a do block in a brace block",
			input:    "for a in [1].map { |v| [v].each do |w| end; v } do last = a end; last",
			expected: 1,
		},
		{
			name:     "a block with do on one line ending in an identifier",
			input:    "s = 0; [5].each do |i| s = i end; s",
			expected: 5,
		},
		{
			name:     "with an empty body",
			input:    "for item in [1, 2]; end",
			expected: []any{1, 2},
		},
		{
			name:     "returns what each returns",
			input:    "r = for x in [1, 2, 3] do y = x * 2 end; [r, x, y]",
			expected: []any{[]any{1, 2, 3}, 3, 6},
		},
		{
			name:     "evaluates the collection once",
			input:    "$calls = 0; def items; $calls += 1; [1, 2]; end; for x in items; end; $calls",
			expected: 1,
		},
		{
			name: "over a class that defines each",
			input: `
				class Pair
					def each
						yield 1
						yield 2
						:done
					end
				end
				seen = []
				result = for item in Pair.new
					seen << item
				end
				[seen, result]
			`,
			expected: []any{[]any{1, 2}, ":done"},
		},
		{
			name:     "over a hash",
			input:    "keys = []; for key in {a: 1, b: 2}; keys << key; end; keys",
			expected: []any{":a", ":b"},
		},
		{
			name:     "overwrites an outer local",
			input:    "x = 9; for x in [4]; end; x",
			expected: 4,
		},
		{
			name:     "defines the variable when the loop runs zero times",
			input:    "for x in []; y = 1; end; [x, y]",
			expected: []any{nil, nil},
		},
		{
			name: "inside a method",
			input: `
				def last_doubled(items)
					for item in items
						doubled = item * 2
					end
					[item, doubled]
				end
				last_doubled([1, 2])
			`,
			expected: []any{2, 4},
		},
		{
			name: "inside a block",
			input: `
				[1].map do |n|
					for item in [n, 5]
						total = item
					end
					[item, total]
				end
			`,
			expected: []any{[]any{5, 5}},
		},
		{
			name: "keeps the locals of a block in its body",
			input: `
				def inner
					:method
				end
				for item in [1]
					[2].each { |n| inner = n }
				end
				inner
			`,
			expected: ":method",
		},
		{
			name: "nested",
			input: `
				pairs = []
				for a in [1, 2]
					for b in [3, 4]
						pairs << [a, b]
					end
				end
				[pairs, a, b]
			`,
			expected: []any{[]any{[]any{1, 3}, []any{1, 4}, []any{2, 3}, []any{2, 4}}, 2, 4},
		},
		{
			name: "rescues an error raised in the body",
			input: `
				def safely
					for item in [1]
						raise "boom"
					end
				rescue => error
					error.message
				end
				safely
			`,
			expected: "boom",
		},
	}
	runVmTests(t, tests)
}
