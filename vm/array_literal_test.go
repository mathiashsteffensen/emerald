package vm

import "testing"

func TestArrayLiteral(t *testing.T) {
	tests := []vmTestCase{
		{"", "[]", []int{}},
		{"", "[1, 2, 3]", []int{1, 2, 3}},
		{"", "[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}
	runVmTests(t, tests)
}
