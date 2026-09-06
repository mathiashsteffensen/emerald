package vm

import "testing"

func TestConditionalSetterEvaluation(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{input: "$hash = {key: 10}; value = (receiver[key] ||= rhs); [value, $receivers, $keys, $rhs]", expected: []any{10, 1, 1, 0}},
		{input: "$hash = {key: false}; value = (receiver[key] &&= rhs); [value, $receivers, $keys, $rhs]", expected: []any{false, 1, 1, 0}},
		{input: "$hash = {}; value = (receiver[key] ||= rhs); [value, $hash[:key], $receivers, $keys, $rhs]", expected: []any{7, 7, 1, 1, 1}},
		{input: "$hash = {key: 10}; value = (receiver[key] &&= rhs); [value, $hash[:key], $receivers, $keys, $rhs]", expected: []any{7, 7, 1, 1, 1}},
	}, `$receivers = 0; $keys = 0; $rhs = 0
	def receiver; $receivers += 1; $hash; end
	def key; $keys += 1; :key; end
	def rhs; $rhs += 1; 7; end`)
}

func TestCompoundSetterAssignment(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "property assignment", input: "item.value += 2; item.value", expected: 12},
		{name: "index assignment", input: "hash = {key: 10}; hash[:key] += 2; hash[:key]", expected: 12},
		{name: "receiver evaluated once", input: `
			$calls = 0
			def receiver; $calls += 1; item; end
			result = (receiver.value *= 2)
			[result, item.value, $calls]
		`, expected: []any{20, 20, 1}},
		{name: "index evaluated once", input: `
			$calls = 0
			def key; $calls += 1; :key; end
			hash = {key: 10}
			hash[key] -= 2
			[hash[:key], $calls]
		`, expected: []any{8, 1}},
	}, `class Item
		def initialize; @value = 10; end
		def value; @value; end
		def value=(n); @value = n; -1; end
	end
	item = Item.new`)
}

func TestSetterAssignmentValue(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "setter returns assigned value", input: `
			class Item
				def value=(n); @value = n; -1; end
				def value; @value; end
			end
			item = Item.new
			result = (item.value = 10)
			[result, item.value]
		`, expected: []any{10, 10}},
		{name: "builtin index setter", input: "hash = {}; hash[:key] = 10", expected: 10},
	})
}

func TestSelf(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "assignment call to self",
			input: `
				def value=(new)
					@value = new
				end

				self.value = 5
				value = 10

				@value
			`,
			expected: 5,
		},
	}

	runVmTests(t, tests)
}
