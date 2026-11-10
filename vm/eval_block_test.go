package vm

import "testing"

func TestRawEvalBlockKeywordLayout(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "zero keywords", input: "item.call(42)", expected: 42},
		{name: "one extra keyword", input: "item.call(42, extra: 1)", expected: 42},
		{name: "two extra keywords", input: "item.call(42, extra: 1, other: 2)", expected: 42},
	}, `item = Proc.new { |value| value }`)
	runVmTests(t, []vmTestCase{
		{name: "declared keywords reordered", input: "pair.call(42, right: 2, left: 1)", expected: []any{42, 1, 2}},
		{name: "declared and extra keywords", input: "pair.call(42, right: 2, extra: 9, left: 1)", expected: []any{42, 1, 2}},
	}, `pair = Proc.new { |value, left:, right:| [value, left, right] }`)
}
