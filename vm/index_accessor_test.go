package vm

import "testing"

func TestIndexMethodDefinitions(t *testing.T) {
	runVmTests(t, []vmTestCase{{
		input: `class Item
			def [](key); @value; end
			def []=(key, n); @value = n; -1; end
		end
		item = Item.new
		result = (item[:key] = 10)
		[result, item[:key]]`,
		expected: []any{10, 10},
	}})
}

func TestIndexAccessor(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "array index accessor",
			input:    "[0,1,2][1]",
			expected: 1,
		},
	}

	runVmTests(t, tests)
}

func TestSpacedLeftBracket(t *testing.T) {
	methods := `
		def count_of(items = [7, 8, 9]); items.size; end
		def pair; [5, 6]; end
		class Box; def count_of(items = [7, 8, 9]); items.size; end; end
		class Failure < StandardError; def [](i); i + 1; end; end
		class Operators
			def ==(a); "eq#{a.size}"; end
			def ===(a); "eqq#{a.size}"; end
			def <=>(a); "cmp#{a.size}"; end
			def <<(a); "app#{a.size}"; end
		end
	`

	tests := []vmTestCase{
		{name: "command argument with a method call", input: "count_of [1, 2].map { |n| n }", expected: 2},
		{name: "index on a call result", input: "pair[1]", expected: 6},
		{name: "receiver call argument", input: "Box.new.count_of [1]", expected: 1},
		{name: "index on a receiver call result", input: "class Box; def items; [5, 6]; end; end; Box.new.items[1]", expected: 6},
		{name: "yield argument", input: "def each_pair; yield [3, 4]; end; each_pair { |a| a[1] }", expected: 4},
		{name: "local", input: "x = [1, 2]; x [0]", expected: 1},
		{name: "method parameter", input: "def first(a); a [0]; end; first([4])", expected: 4},
		{name: "block parameter", input: "[[4]].map { |a| a [0] }", expected: []any{4}},
		{name: "enclosing local in a block", input: "x = [4]; [0].map { |i| x [i] }", expected: []any{4}},
		{name: "for variable", input: "r = 0; for a in [[4]]; r = a [0]; end; r", expected: 4},
		{name: "rescue variable", input: "def fail_and_index; raise Failure; rescue => e; e [1]; end; fail_and_index", expected: 2},
		{name: "non-ASCII local", input: "café = [4]; café [0]", expected: 4},
		{name: "local of another method is a call", input: "def a; pair = [1]; end; def b; pair [0]; end; def pair(x); x.size; end; b", expected: 1},
		{name: "== argument", input: "Operators.new.== [1]", expected: "eq1"},
		{name: "=== argument", input: "Operators.new.=== [1, 2]", expected: "eqq2"},
		{name: "<=> argument", input: "Operators.new.<=> [1, 2, 3]", expected: "cmp3"},
		{name: "<< argument", input: "Operators.new.<< [1]", expected: "app1"},
	}

	for i := range tests {
		tests[i].input = methods + tests[i].input
	}

	runVmTests(t, tests)
}
