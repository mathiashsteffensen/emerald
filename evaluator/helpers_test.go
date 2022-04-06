package evaluator

import (
	"emerald/lexer"
	"emerald/object"
	"emerald/parser"
	"testing"
)

func testEval(input string) object.EmeraldValue {
	l := lexer.New(lexer.NewInput("test.rb", input))
	p := parser.New(l)
	program := p.ParseAST()
	return Eval(object.Object, program, object.NewEnvironment())
}

func testNullObject(t *testing.T, obj object.EmeraldValue) bool {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.EmeraldValue, expected bool) bool {
	class := obj.ParentClass().(*object.Class).Name
	if class != "TrueClass" && class != "FalseClass" {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	actual := class == "TrueClass"

	if actual != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", actual, expected)
		return false
	}

	return true
}

func testIntegerObject(t *testing.T, obj object.EmeraldValue, expected int64) bool {
	result, ok := obj.(*object.IntegerInstance)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testSymbolObject(t *testing.T, obj object.EmeraldValue, expected string) bool {
	result, ok := obj.(*object.SymbolInstance)
	if !ok {
		t.Errorf("object is not Symbol. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}
	return true
}
