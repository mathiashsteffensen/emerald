package vm

import "testing"

func TestFloatArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"adding", "1.5 + 2.7", 4.2},
		{"subtracting", "1.0 - 2.8 == 1.8", false},
		{"multiplying", "12.5 * 2", 25.0},
		{"dividing", "2.5_0 / 0.5", 5.0},
		{"grouped expression", "5.4 * (2 + 10) == 64.8", false},
		{"negating", "-5.4", -5.4},
		{"negating and adding", "-50.0 + 100 + -50", 0.0},
		{"mixed", "(5.0 + 10 * 2 + 15 / 3) * 2 + -10", 50.0},
	}
	runVmTests(t, tests)
}
