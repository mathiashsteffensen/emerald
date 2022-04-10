package vm

import "testing"

func TestBooleanLiteral(t *testing.T) {
	tests := []vmTestCase{
		{"true", "true", true},
		{"false", "false", false},
	}
	runVmTests(t, tests)
}
