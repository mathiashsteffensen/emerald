package vm

import "testing"

func TestGlobalAssignment(t *testing.T) {
	tests := []vmTestCase{
		{"", "one = 1", 1},
		{"", "one = 1; one", 1},
		{"", "one = 1; two = 2; one + two", 3},
		{"", "one = 1; two = one + one; one + two", 3},
	}
	runVmTests(t, tests)
}
