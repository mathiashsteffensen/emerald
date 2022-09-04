package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/heap"
	"emerald/object"
	"emerald/parser"
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"strings"
	"testing"
)

type vmTestCase struct {
	name     string
	input    string
	expected any
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core.Object.ResetForSpec()

			program := parse(tt.input)
			comp := compiler.New()

			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			vm := New("test", comp.Bytecode())
			safeRun(t, vm)

			stackElem := safePop(t, vm)
			testExpectedObject(t, tt.expected, stackElem)

			if vm.currentFiber().sp != 0 {
				if str, ok := tt.expected.(string); ok && strings.HasPrefix(str, "error:") {
					return
				}

				t.Errorf("stack pointer was not reset after running test, this indicates a memory leak in the VM, was %d", vm.currentFiber().sp)
			}

			heap.Reset()
		})
	}
}

func safeRun(t *testing.T, vm *VM) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("vm error: %s", r)
		}
	}()

	vm.Run()
}

func safePop(t *testing.T, vm *VM) object.EmeraldValue {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Failed to to get last popped stack element with err %s", r)
		}
	}()

	return vm.LastPoppedStackElem()
}

func testExpectedObject(
	t *testing.T,
	expected any,
	actual object.EmeraldValue,
) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case float64:
		err := testFloatObject(expected, actual)
		if err != nil {
			t.Errorf("testFloatObject failed: %s", err)
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
		} else {
			if strings.HasPrefix(expected, "class:") {
				err := testClassObject(expected[6:], actual)
				if err != nil {
					t.Errorf("testClassObject failed: %s", err)
				}
			} else {
				if strings.HasPrefix(expected, "instance:") {
					err := testInstanceObject(expected[9:], actual)
					if err != nil {
						t.Errorf("testInstanceObject failed: %s", err)
					}
				} else {
					if strings.HasPrefix(expected, "error:") {
						err := testErrorObject(expected[6:], actual)
						if err != nil {
							t.Errorf("testErrorObject failed: %s", err)
						}
					} else {
						err := testStringObject(expected, actual)
						if err != nil {
							t.Errorf("testStringObject failed: %s", err)
						}
					}
				}
			}
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

func parse(input string) *ast.AST {
	l := lexer.New(lexer.NewInput("test.rb", input))
	p := parser.New(l)
	return p.ParseAST()
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

	if len(hash.Values) != len(expected) {
		return fmt.Errorf("hash has wrong number of Pairs. want=%d, got=%d", len(expected), len(hash.Values))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hash.Values[expectedKey.HashKey()]
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
		return fmt.Errorf("object is not IntegerInstance. got=%s", actual.Inspect())
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testFloatObject(expected float64, actual object.EmeraldValue) error {
	result, ok := actual.(*core.FloatInstance)
	if !ok {
		return fmt.Errorf("object is not FloatInstance. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
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
	actualClass, ok := actual.(*object.Class)
	if !ok {
		return fmt.Errorf("expected class got=%s", actual.Inspect())
	}

	if expected != actualClass.Name {
		return fmt.Errorf("expectedClass was expected to be %s, got=%s", expected, actualClass.Name)
	}

	return nil
}

func testInstanceObject(expected string, actual object.EmeraldValue) error {
	actualInstance, ok := actual.(*object.Instance)
	if !ok {
		return fmt.Errorf("expected instance got=%#v", actual)
	}

	class := actualInstance.Class().Super().(*object.Class)

	if class.Name != expected {
		return fmt.Errorf("expected instance to be instance of %s, but is instance of %s", expected, class.Name)
	}

	return nil
}

func testErrorObject(expected string, actual object.EmeraldValue) error {
	actualInstance, ok := actual.(object.EmeraldError)
	if !ok {
		return fmt.Errorf("expected StandardErrorInstance got=%T", actual)
	}

	if actualInstance.Message() != expected {
		return fmt.Errorf("expected error to have message %s, got=%s", expected, actualInstance.Message())
	}

	return nil
}
