package vm

import "testing"

func TestBooleanExpression(t *testing.T) {
	tests := []vmTestCase{
		{"true", "true", true},
		{"false", "false", false},
		{"", "1 < 2", true},
		{"", "1 > 2", false},
		{"", "1 < 1", false},
		{"", "1 > 1", false},
		{"", "1 == 1", true},
		{"", "1 != 1", false},
		{"", "1 == 2", false},
		{"", "1 != 2", true},
		{"", "true == true", true},
		{"", "false == false", true},
		{"", "true == false", false},
		{"", "true != false", true},
		{"", "false != true", true},
		{"", "(1 < 2) == true", true},
		{"", "(1 < 2) == false", false},
		{"", "(1 > 2) == true", false},
		{"", "(1 > 2) == false", true},
		{"", "!true", false},
		{"", "!false", true},
		{"", "!5", false},
		{"", "!!true", true},
		{"", "!!false", false},
		{"", "!!5", true},
		{"", "!(if false; 5; end)", true},
	}
	runVmTests(t, tests)
}
