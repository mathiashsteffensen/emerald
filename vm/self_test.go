package vm

import "testing"

func TestRescuedSetterPreservesLocals(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "missing setter", input: `
			def f
				n = 42
				Object.new.nope = 7
			rescue NoMethodError
				n
			end
			f
		`, expected: 42},
		{name: "wrong setter arity", input: `
			class Box; def value=(a, b); a; end; end
			def f
				n = 42
				Box.new.value = 7
			rescue ArgumentError
				n
			end
			f
		`, expected: 42},
		{name: "setter raises to caller", input: `
			class Box; def value=(x); raise "boom"; end; end
			def f
				n = 42
				Box.new.value = 7
			rescue StandardError
				n
			end
			f
		`, expected: 42},
		{name: "setter rescues internally and returns normally", input: `
			class Box
				def value=(x); raise "boom"; rescue StandardError; -1; end
			end
			Box.new.value = 7
		`, expected: 7},
	})
}

func TestConditionalSetterForwardsBlock(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{input: "value = (box.value { 42 } ||= 7); [value, $gets, $sets, $yielded]", expected: []any{42, 1, 0, nil}},
		{input: "value = (box.value { nil } ||= 7); [value, $gets, $sets, $yielded]", expected: []any{7, 1, 1, nil}},
		{input: "value = (box.value { false } &&= 7); [value, $gets, $sets, $yielded]", expected: []any{false, 1, 0, nil}},
		{input: "value = (box.value { 42 } &&= 7); [value, $gets, $sets, $yielded]", expected: []any{7, 1, 1, 42}},
	}, `class Box
		def value; $gets += 1; yield; end
		def value=(x); $sets += 1; $yielded = yield; -1; end
	end
	box = Box.new; $gets = 0; $sets = 0; $yielded = nil`)
}

func TestConditionalSetterForwardsKeywords(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{input: "value = (box.value(5, key: key(42)) ||= 7); [value, $gets, $keys, $received]", expected: []any{42, 1, 1, nil}},
		{input: "value = (box.value(5, key: key(nil)) ||= 7); [value, $gets, $keys, $received]", expected: []any{7, 1, 1, []any{5, 7, nil}}},
		{input: "value = (box.value(5, key: key(false)) &&= 7); [value, $gets, $keys, $received]", expected: []any{false, 1, 1, nil}},
		{input: "value = (box.value(5, key: key(42)) &&= 7); [value, $gets, $keys, $received]", expected: []any{7, 1, 1, []any{5, 7, 42}}},
		{input: "value = (box.value(5, key: key(42)) += 7); [value, $gets, $keys, $received]", expected: []any{49, 1, 1, []any{5, 49, 42}}},
	}, `class Box
		def value(index, key:); $gets += 1; key; end
		def value=(index, x, key:); $received = [index, x, key]; -1; end
	end
	def key(value); $keys += 1; value; end
	box = Box.new; $gets = 0; $keys = 0; $received = nil`)
}

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

func TestSetterAssignmentValueWithKeywords(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "integer RHS", input: "box.value(key: 42) = 7", expected: 7},
		{name: "nil RHS", input: "box.value(key: 42) = nil", expected: nil},
		{name: "false RHS", input: "box.value(key: 42) = false", expected: false},
		{name: "keyword forwarding", input: "box.value(key: 42) = 7; $received", expected: []any{7, 42}},
	}, `class Box
		def value=(x, key:); $received = [x, key]; -1; end
	end
	box = Box.new`)
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
