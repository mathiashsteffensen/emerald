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
