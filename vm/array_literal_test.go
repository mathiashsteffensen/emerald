package vm

import "testing"

func TestArrayLiteral(t *testing.T) {
	tests := []vmTestCase{
		{"", "[]", []any{}},
		{"", "[1, 2, 3]", []any{1, 2, 3}},
		{"", "[1 + 2, 3 * 4, 5 + 6]", []any{3, 12, 11}},
	}
	runVmTests(t, tests)
}
