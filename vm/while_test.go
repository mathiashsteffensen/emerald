package vm

import "testing"

func TestWhile(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				new_items = []
				items = [1, 4, 9]
				while item = items.pop
					new_items.push(item * 2)
				end
				new_items
			`,
			expected: []any{18, 8, 2},
		},
		{
			name: "as modifier",
			input: `
				new_items = []
				items = [1, 4, 9]
				new_items.push(2) while item = items.pop
				new_items
			`,
			expected: []any{2, 2, 2},
		},
	}
	runVmTests(t, tests)
}
