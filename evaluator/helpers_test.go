package evaluator

import (
	"emerald/lexer"
	"emerald/object"
	"emerald/parser"
	"testing"
)

func testEval(input string, envs ...object.Environment) object.EmeraldValue {
	l := lexer.New(lexer.NewInput("test.rb", input))
	p := parser.New(l)
	program := p.ParseAST()
	env := object.NewEnvironment()

	if len(envs) != 0 {
		env = envs[0]
	}

	return Eval(object.ExecutionContext{Target: object.Object, IsStatic: true}, program, env)
}

func testObjectValue(t *testing.T, obj object.EmeraldValue, expected any) bool {
	t.Helper()

	if expected == nil {
		return testNullObject(t, obj)
	}

	switch exp := expected.(type) {
	case bool:
		return testBooleanObject(t, obj, exp)
	case int:
		return testIntegerObject(t, obj, int64(exp))
	case int64:
		return testIntegerObject(t, obj, exp)
	case string:
		return testSymbolObject(t, obj, exp)
	}

	t.Fatalf("don't know how to test %T values", expected)
	return false
}

func testNullObject(t *testing.T, obj object.EmeraldValue) bool {
	t.Helper()

	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.EmeraldValue, expected bool) bool {
	t.Helper()

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
	t.Helper()

	result, ok := obj.(*object.IntegerInstance)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v) (%s)", obj, obj, obj.Inspect())
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
	t.Helper()

	result, ok := obj.(*object.SymbolInstance)
	if !ok {
		t.Fatalf("object is not Symbol. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Fatalf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}
	return true
}
