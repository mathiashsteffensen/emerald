package vm

import "testing"

func TestUntil(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "runs while the condition is falsy",
			input: `
				items = []
				i = 0
				until i == 3
					i += 1
					items.push(i)
				end
				items
			`,
			expected: []any{1, 2, 3},
		},
		{
			name: "with a semicolon after the condition",
			input: `
				done = nil
				until done; done = true; end
				done
			`,
			expected: true,
		},
		{
			name: "skips the body when the condition is truthy at the start",
			input: `
				ran = false
				until true
					ran = true
				end
				ran
			`,
			expected: false,
		},
		{
			name: "is nil",
			input: `
				i = 0
				until i == 2
					i += 1
				end
			`,
			expected: nil,
		},
		{
			name: "as modifier",
			input: `
				i = 0
				i += 1 until i >= 3
				i
			`,
			expected: 3,
		},
		{
			name: "as modifier skips the body when the condition is truthy at the start",
			input: `
				i = 5
				i += 1 until i >= 3
				i
			`,
			expected: 5,
		},
		{
			name: "break with a value",
			input: `
				i = 0
				until false
					i += 1
					break i * 10 if i == 3
				end
			`,
			expected: 30,
		},
		{
			name:     "bare break",
			input:    `until false; break; end`,
			expected: nil,
		},
		{
			name: "next skips to the condition",
			input: `
				odds = []
				i = 0
				until i == 5
					i += 1
					next if i == 2 || i == 4
					odds.push(i)
				end
				odds
			`,
			expected: []any{1, 3, 5},
		},
		{
			name: "nested in a while loop and a block",
			input: `
				pairs = []
				i = 0
				while i < 2
					i += 1
					j = 0
					until j == 2
						j += 1
						[i].each { |x| pairs.push([x, j]) }
						break if j == 1 && i == 2
					end
				end
				pairs
			`,
			expected: []any{[]any{1, 1}, []any{1, 2}, []any{2, 1}},
		},
	}
	runVmTests(t, tests)
}

// A `while` or `until` modifier on `break`, `next` or `return` is a loop of its own around that statement
func TestLoopModifierOnJumpStatements(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "break until a truthy condition never runs",
			input: `
				items = []
				i = 0
				while i < 3
					i += 1
					break until true
					items.push(i)
				end
				items
			`,
			expected: []any{1, 2, 3},
		},
		{
			name: "break until a falsy condition leaves only the modifier loop",
			input: `
				items = []
				i = 0
				while i < 3
					i += 1
					break until false
					items.push(i)
				end
				items
			`,
			expected: []any{1, 2, 3},
		},
		{
			name: "break while a truthy condition leaves only the modifier loop",
			input: `
				items = []
				i = 0
				while i < 3
					i += 1
					break 7 while true
					items.push(i)
				end
				items
			`,
			expected: []any{1, 2, 3},
		},
		{
			name: "next until repeats the modifier loop",
			input: `
				items = []
				i = 0
				while i < 3
					i += 1
					n = 0
					next until (n += 1) > 2
					items.push([i, n])
				end
				items
			`,
			expected: []any{[]any{1, 3}, []any{2, 3}, []any{3, 3}},
		},
		{
			name: "next while a falsy condition never runs",
			input: `
				items = []
				i = 0
				until i == 3
					i += 1
					next while false
					items.push(i)
				end
				items
			`,
			expected: []any{1, 2, 3},
		},
		{
			name:     "next until a truthy condition in a block never runs",
			input:    `[1, 2].map { |v| next v * 10 until true; v }`,
			expected: []any{1, 2},
		},
		{
			name: "return until a falsy condition returns",
			input: `
				def until_return
					return 1 until false
					3
				end
				until_return
			`,
			expected: 1,
		},
		{
			name: "return while a falsy condition falls through",
			input: `
				def while_return
					return 2 while false
					3
				end
				while_return
			`,
			expected: 3,
		},
	}
	runVmTests(t, tests)
}

// Modifiers after `break`, `next` or `return` apply left to right: `return 1 if a unless b` is `(return 1 if a) unless b`
func TestChainedModifiersOnJumpStatements(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "return if unless",
			input: `
				def chained_return(a, b)
					return 1 if a unless b
					2
				end
				[chained_return(true, true), chained_return(true, false), chained_return(false, true), chained_return(false, false)]
			`,
			expected: []any{2, 1, 2, 2},
		},
		{
			name: "bare return if unless",
			input: `
				def chained_bare_return(a, b)
					return if a unless b
					2
				end
				[chained_bare_return(true, false), chained_bare_return(false, false)]
			`,
			expected: []any{nil, 2},
		},
		{
			name: "next if while",
			input: `
				items = []
				[1, 2, 3].each do |i|
					k = 0
					next if i == 2 while (k += 1) < 3
					items.push(i)
				end
				items
			`,
			expected: []any{1, 2, 3},
		},
		{
			name: "break if unless",
			input: `
				def chained_break(a, b)
					[1, 2, 3].each do |i|
						break i * 10 if a unless b
					end
				end
				[chained_break(true, true), chained_break(true, false), chained_break(false, true), chained_break(false, false)]
			`,
			expected: []any{[]any{1, 2, 3}, 10, []any{1, 2, 3}, []any{1, 2, 3}},
		},
		{
			name: "break if unless in a while loop",
			input: `
				x = 0
				while true
					x += 1
					break x if x > 2 unless x == 3
				end
			`,
			expected: 4,
		},
	}
	runVmTests(t, tests)
}
