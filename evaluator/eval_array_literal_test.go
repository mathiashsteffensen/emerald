package evaluator

import (
	"emerald/object"
	"testing"
)

func TestArrayLiteralParsing(t *testing.T) {
	input := `identifier = 1
	[0, identifier, 2]`

	expected := []int64{0, 1, 2}

	evaluated := testEval(input)

	arr, ok := evaluated.(*object.ArrayInstance)
	if !ok {
		t.Fatalf("evaluated was not ArrayInstance, got=%T (%+v)", evaluated, evaluated)
	}

	for i, exp := range expected {
		testIntegerObject(t, arr.Value[i], exp)
	}
}
