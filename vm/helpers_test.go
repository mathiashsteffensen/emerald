package vm

import (
	"emerald/bytecode"
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"fmt"
	"math"
	"strings"
	"testing"
)

type vmTestCase struct {
	name     string
	input    string
	expected any
}

func runVmTests(t *testing.T, tests []vmTestCase, setupScripts ...string) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputs := make([]string, len(setupScripts))
			copy(inputs, setupScripts)
			inputs = append(inputs, tt.input)

			rt := core.NewRuntime()
			rt.Init()

			rt.CompileBlock = func(fileName string, content string) *bytecode.Bytecode {
				return compiler.Compile(fileName, content, rt)
			}

			bc := compiler.Compile("test", strings.Join(inputs, "\n"), rt)

			vm := New("test", bc, rt)
			vm.Run()

			ensureNoExceptionUnlessExpected(t, tt.expected, rt)

			stackElem := safePop(t, vm)
			testExpectedObject(t, tt.expected, stackElem, rt)

			if vm.currentFiber().sp != 0 {
				if str, ok := tt.expected.(string); ok && strings.HasPrefix(str, "error:") {
					return
				}

				t.Errorf("stack pointer was not reset after running test, this indicates a memory leak in the VM, was %d", vm.currentFiber().sp)
			}
		})
	}
}

func safePop(t *testing.T, vm *VM) object.EmeraldValue {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Failed to to get last popped stack element with err %s", r)
		}
	}()

	return vm.LastPoppedStackElem()
}

func ensureNoExceptionUnlessExpected(t *testing.T, expected any, rt *core.Runtime) {
	if expectedString, ok := expected.(string); ok && strings.HasPrefix(expectedString, "error:") {
		return
	}

	exception := rt.Heap.GetGlobalVariableString("$!")

	if !exception.IsNil() {
		t.Fatalf("Unexpected uncaught exception %s: %s", exception.Inspect(), exception.Heap.(object.EmeraldError).Message())
	}
}

func testExpectedObject(
	t *testing.T,
	expected any,
	actual object.EmeraldValue,
	rt *core.Runtime,
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
		err := testBooleanObject(expected, actual, rt)
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
						err := testErrorObject(expected[6:], rt.Heap.GetGlobalVariableString("$!"))
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
		err := testArrayObject(t, expected, actual, rt)
		if err != nil {
			t.Errorf("testArrayObject failed: %s", err)
		}
	case map[object.EmeraldValue]any:
		err := testHashObject(t, expected, actual, rt)
		if err != nil {
			t.Errorf("testHashObject failed: %s", err)
		}
	case nil:
		if !actual.IsNil() {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}
	}
}

func testArrayObject(t *testing.T, expected []any, actual object.EmeraldValue, rt *core.Runtime) error {
	array, ok := actual.Heap.(*core.ArrayInstance)
	if !ok {
		return fmt.Errorf("object not Array: %T (%+v)", actual, actual)
	}

	if len(array.Value) != len(expected) {
		return fmt.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Value))
	}

	for i, expectedElem := range expected {
		testExpectedObject(t, expectedElem, array.Value[i], rt)
	}

	return nil
}

func testHashObject(t *testing.T, expected map[object.EmeraldValue]any, actual object.EmeraldValue, rt *core.Runtime) error {
	hash, ok := actual.Heap.(*core.HashInstance)
	if !ok {
		return fmt.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
	}

	if hash.Values.Len() != len(expected) {
		return fmt.Errorf("hash has wrong number of Pairs. want=%d, got=%d", len(expected), hash.Values.Len())
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hash.Values.Get(expectedKey.HashKey())
		if !ok {
			return fmt.Errorf("no pair for given key in Pairs")
		}

		testExpectedObject(t, expectedValue, pair, rt)
	}

	return nil
}

func testIntegerObject(expected int64, actual object.EmeraldValue) error {
	if !actual.Is(object.INTEGER_VALUE) {
		return fmt.Errorf("object is not Integer. got=%s", actual.Inspect())
	}
	if int64(actual.Num) != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", int64(actual.Num), expected)
	}
	return nil
}

func testFloatObject(expected float64, actual object.EmeraldValue) error {
	if !actual.Is(object.FLOAT_VALUE) {
		return fmt.Errorf("object is not Float. got=%T (%+v)", actual, actual)
	}

	val := math.Float64frombits(actual.Num)
	if val != expected {
		return fmt.Errorf("object has wrong value. got=%f, want=%f", val, expected)
	}
	return nil
}

func testBooleanObject(expected bool, actual object.EmeraldValue, rt *core.Runtime) error {
	if !actual.Is(object.TRUE_VALUE) && !actual.Is(object.FALSE_VALUE) {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if actual.Is(object.TRUE_VALUE) != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", actual.Is(object.TRUE_VALUE), expected)
	}
	return nil
}

func testStringObject(expected string, actual object.EmeraldValue) error {
	result, ok := actual.Heap.(*core.StringInstance)
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
	result, ok := actual.Heap.(*core.SymbolInstance)
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
	actualClass, ok := actual.Heap.(*object.Class)
	if !ok {
		return fmt.Errorf("expected class got=%s", actual.Inspect())
	}

	if expected != actualClass.Name {
		return fmt.Errorf("expectedClass was expected to be %s, got=%s", expected, actualClass.Name)
	}

	return nil
}

func testInstanceObject(expected string, actual object.EmeraldValue) error {
	class := object.RealClass(actual).Heap.(*object.Class)

	if class.Name != expected {
		return fmt.Errorf("expected instance to be instance of %s, but is instance of %s", expected, class.Name)
	}

	return nil
}

func testErrorObject(expected string, actual object.EmeraldValue) error {
	split := strings.Split(expected, ":")
	className := split[0]
	msg := strings.Join(split[1:], ":")

	emeraldError, ok := actual.Heap.(object.EmeraldError)
	if !ok {
		return fmt.Errorf("object was not EmeraldError, got=%T", actual)
	}

	if emeraldError.ClassName() != className {
		return fmt.Errorf("unexpected error class \nwant=%s\ngot=%s", className, emeraldError.ClassName())
	}

	if emeraldError.Message() != msg {
		return fmt.Errorf("unexpected error msg \nwant=%s\ngot=%s", msg, emeraldError.Message())
	}

	return nil
}
