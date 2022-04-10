package vm

import "testing"

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"literal", "1", 1},
		{"adding", "1 + 2", 3},
		{"subtracting", "1 - 2", -1},
		{"multiplying", "1 * 2", 2},
		{"dividing", "4 / 2", 2},
		{"mixed", "50 / 2 * 2 + 10 - 5", 55},
		{"grouped expression", "5 * (2 + 10)", 60},
	}
	runVmTests(t, tests)
}
