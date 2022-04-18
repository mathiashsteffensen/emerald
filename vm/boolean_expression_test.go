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
		{"", "1 >= 1", true},
		{"", "1 >= 2", false},
		{"", "1 <= 1", true},
		{"", "2 <= 1", false},
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
		{"false expression == true", "(1 > 2) == true", false},
		{"false expression == false", "(1 > 2) == false", true},
		{"negate true", "!true", false},
		{"negate false", "!false", true},
		{"negate integer", "!5", false},
		{"boolean force true", "!!true", true},
		{"boolean force false", "!!false", false},
		{"boolean force integer", "!!5", true},
		{"negate if expression resolving to nil", "!(if false; 5; end)", true},
		{"boolean and evaluating to true", "true && 15", 15},
		{"boolean and evaluating to false", "false && 15", false},
		{"boolean or evaluating to true", "true || 15", true},
		{"boolean or evaluating to false", "false || 15", 15},
	}
	runVmTests(t, tests)
}
