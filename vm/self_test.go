package vm

import "testing"

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
