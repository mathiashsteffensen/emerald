package vm

import "testing"

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"adding", "1 + 2", 3},
		{"subtracting", "1 - 2", -1},
		{"multiplying", "1 * 2", 2},
		{"dividing", "4 / 2", 2},
		{"grouped expression", "5 * (2 + 10)", 60},
		{"negating", "-5", -5},
		{"negating and adding", "-50 + 100 + -50", 0},
		{"mixed", "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVmTests(t, tests)
}
