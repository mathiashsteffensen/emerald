package vm

import "testing"

func TestRawEvalBlockKeywordLayout(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "zero keywords", input: "Item.new(42).value", expected: 42},
		{name: "one extra keyword", input: "Item.new(42, extra: 1).value", expected: 42},
		{name: "two extra keywords", input: "Item.new(42, extra: 1, other: 2).value", expected: 42},
	}, `class Item
		define_method(:initialize) { |value| @value = value }
		def value; @value; end
	end`)
	runVmTests(t, []vmTestCase{
		{name: "declared keywords reordered", input: "Pair.new(42, right: 2, left: 1).values", expected: []any{42, 1, 2}},
		{name: "declared and extra keywords", input: "Pair.new(42, right: 2, extra: 9, left: 1).values", expected: []any{42, 1, 2}},
	}, `class Pair
		define_method(:initialize) { |value, left:, right:| @values = [value, left, right] }
		def values; @values; end
	end`)
}
