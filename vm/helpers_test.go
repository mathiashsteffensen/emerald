package vm

import (
	"emerald/ast"
	"emerald/compiler"
	"emerald/lexer"
	"emerald/object"
	"emerald/parser"
	"fmt"
	"testing"
)

type vmTestCase struct {
	name     string
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parse(tt.input)
			comp := compiler.New()

			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			vm := New(comp.Bytecode())

			err = vm.Run()
			if err != nil {
				t.Fatalf("vm error: %s", err)
			}

			stackElem := vm.LastPoppedStackElem()
			testExpectedObject(t, tt.expected, stackElem)
		})
	}
}

func testExpectedObject(
	t *testing.T,
	expected interface{},
	actual object.EmeraldValue,
) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []any:
		testArrayObject(t, expected, actual)
	case nil:
		if actual != object.NULL {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}
	}
}

func parse(input string) *ast.AST {
	l := lexer.New(lexer.NewInput("test.rb", input))
	p := parser.New(l)
	return p.ParseAST()
}

func testArrayObject(t *testing.T, expected []any, actual object.EmeraldValue) {
	array, ok := actual.(*object.ArrayInstance)
	if !ok {
		t.Errorf("object not Array: %T (%+v)", actual, actual)
		return
	}

	if len(array.Value) != len(expected) {
		t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Value))
		return
	}

	for i, expectedElem := range expected {
		testExpectedObject(t, expectedElem, array.Value[i])
	}
}

func testIntegerObject(expected int64, actual object.EmeraldValue) error {
	result, ok := actual.(*object.IntegerInstance)
	if !ok {
		return fmt.Errorf("object is not IntegerInstance. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testBooleanObject(expected bool, actual object.EmeraldValue) error {
	if actual != object.TRUE && actual != object.FALSE {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if (actual == object.TRUE) != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", actual == object.TRUE, expected)
	}
	return nil
}

func testStringObject(expected string, actual object.EmeraldValue) error {
	result, ok := actual.(*object.StringInstance)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)",
			actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
	}
	return nil
}
