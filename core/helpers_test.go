package core_test

import (
	"emerald/ast"
	"emerald/compiler"
	"emerald/core"
	"emerald/lexer"
	"emerald/object"
	"emerald/parser"
	"emerald/vm"
	"fmt"
	"strings"
	"testing"
)

type coreTestCase struct {
	name     string
	input    string
	expected any
}

func runCoreTests(t *testing.T, tests []coreTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, class := range object.Classes {
				class.ResetForSpec()
				class.Class().(*object.SingletonClass).ResetForSpec()
			}

			for _, module := range object.Modules {
				module.ResetForSpec()
				module.Class().(*object.SingletonClass).ResetForSpec()
			}

			program := parse(t, tt.input)
			comp := compiler.New(compiler.WithBuiltIns())

			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			machine := vm.New(comp.Bytecode())

			err = machine.Run()
			if err != nil {
				t.Fatalf("vm error: %s", err)
			}

			stackElem := machine.LastPoppedStackElem()
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
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		if strings.HasPrefix(expected, ":") {
			err := testSymbolObject(expected[1:], actual)
			if err != nil {
				t.Errorf("testSymbolObject failed: %s", err)
			}
			return
		}

		if strings.HasPrefix(expected, "class:") {
			err := testClassObject(expected[6:], actual)
			if err != nil {
				t.Errorf("testClassObject failed: %s", err)
			}
			return
		}

		if strings.HasPrefix(expected, "instance:") {
			err := testInstanceObject(expected[9:], actual)
			if err != nil {
				t.Errorf("testInstanceObject failed: %s", err)
			}
			return
		}

		if strings.HasPrefix(expected, "error:") {
			err := testErrorObject(expected[6:], actual)
			if err != nil {
				t.Errorf("testErrorObject failed: %s", err)
			}
			return
		}

		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []any:
		err := testArrayObject(t, expected, actual)
		if err != nil {
			t.Errorf("testArrayObject failed: %s", err)
		}
	case map[object.EmeraldValue]any:
		err := testHashObject(t, expected, actual)
		if err != nil {
			t.Errorf("testHashObject failed: %s", err)
		}
	case nil:
		if actual != core.NULL {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}
	}
}

func parse(t *testing.T, input string) *ast.AST {
	t.Helper()

	l := lexer.New(lexer.NewInput("test.rb", input))
	p := parser.New(l)
	program := p.ParseAST()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			t.Errorf("parser error: %s\n", err)
		}
	}

	return program
}

func testArrayObject(t *testing.T, expected []any, actual object.EmeraldValue) error {
	array, ok := actual.(*core.ArrayInstance)
	if !ok {
		return fmt.Errorf("object not Array: %T (%+v)", actual, actual)
	}

	if len(array.Value) != len(expected) {
		return fmt.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Value))
	}

	for i, expectedElem := range expected {
		testExpectedObject(t, expectedElem, array.Value[i])
	}

	return nil
}

func testHashObject(t *testing.T, expected map[object.EmeraldValue]any, actual object.EmeraldValue) error {
	hash, ok := actual.(*core.HashInstance)
	if !ok {
		return fmt.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
	}

	if len(hash.Value) != len(expected) {
		return fmt.Errorf("hash has wrong number of Pairs. want=%d, got=%d", len(expected), len(hash.Value))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hash.Value[expectedKey]
		if !ok {
			return fmt.Errorf("no pair for given key in Pairs")
		}

		testExpectedObject(t, expectedValue, pair)
	}

	return nil
}

func testIntegerObject(expected int64, actual object.EmeraldValue) error {
	result, ok := actual.(*core.IntegerInstance)
	if !ok {
		return fmt.Errorf("object is not IntegerInstance. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testBooleanObject(expected bool, actual object.EmeraldValue) error {
	if actual != core.TRUE && actual != core.FALSE {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if (actual == core.TRUE) != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", actual == core.TRUE, expected)
	}
	return nil
}

func testStringObject(expected string, actual object.EmeraldValue) error {
	result, ok := actual.(*core.StringInstance)
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

func testSymbolObject(expected string, actual object.EmeraldValue) error {
	result, ok := actual.(*core.SymbolInstance)
	if !ok {
		return fmt.Errorf("object is not Symbol. got=%T (%+v)",
			actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
	}
	return nil
}

func testClassObject(expected string, actual object.EmeraldValue) error {
	expectedClass, ok := object.GetClassByName(expected)
	if !ok {
		return fmt.Errorf("undefined class %s", expected)
	}

	actualClass, ok := actual.(*object.Class)
	if !ok {
		return fmt.Errorf("expected class got=%T", actual)
	}

	if expectedClass.Name != actualClass.Name {
		return fmt.Errorf("expectedClass was expected to be %s, got=%s", expectedClass.Name, actualClass.Name)
	}

	return nil
}

func testInstanceObject(expected string, actual object.EmeraldValue) error {
	actualInstance, ok := actual.(*object.Instance)
	if !ok {
		return fmt.Errorf("expected instance got=%T", actual)
	}

	class := actualInstance.Class().Super().(*object.Class)

	if class.Name != expected {
		return fmt.Errorf("expected instance to be instance of %s, but is instance of %s", expected, class.Name)
	}

	return nil
}

func testErrorObject(expected string, actual object.EmeraldValue) error {
	split := strings.Split(expected, ":")
	className := split[0]
	msg := split[1]

	emeraldError, ok := actual.(object.EmeraldError)
	if !ok {
		return fmt.Errorf("object was not EmeraldError, got=%#v", actual)
	}

	if emeraldError.ClassName() != className {
		return fmt.Errorf("unexpected error class \nwant=%s\ngot=%s", className, emeraldError.ClassName())
	}

	if emeraldError.Message() != msg {
		return fmt.Errorf("unexpected error msg \nwant=%s\ngot=%s", msg, emeraldError.Message())
	}

	return nil
}
