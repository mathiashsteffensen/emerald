package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/object"
	"emerald/parser"
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"math"
	"strings"
	"testing"
)

type compilerTestCase struct {
	name                 string
	input                string
	expectedConstants    []any
	expectedInstructions []bytecode.Instructions
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, program := parse(tt.input)

			rt := core.NewRuntime()
			rt.Init()

			compiler := New(l, rt)

			compiler.Compile(program)

			bytecode := compiler.Bytecode()

			err := testInstructions(tt.expectedInstructions, bytecode.Instructions)
			if err != nil {
				t.Errorf("testInstructions failed: %s", err)
			}

			err = testConstants(tt.expectedConstants, rt.Heap.ConstantPool)

			if err != nil {
				t.Errorf("testConstants failed: %s", err)
			}
		})
	}
}

func parse(input string) (*lexer.Lexer, *ast.AST) {
	l := lexer.New(lexer.NewInput("test.rb", input))
	p := parser.New(l)
	return l, p.ParseAST()
}

func testInstructions(
	expected []bytecode.Instructions,
	actual bytecode.Instructions,
) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot=%q", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot=%q", i, concatted, actual)
		}
	}

	return nil
}

func concatInstructions(s []bytecode.Instructions) bytecode.Instructions {
	out := bytecode.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testConstants(
	expected []any,
	actual []object.EmeraldValue,
) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d",
			len(actual), len(expected))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s",
					i, err)
			}
		case float64:
			err := testFloatObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testFloatObject failed: %s", i, err)
			}
		case string:
			if strings.HasPrefix(constant, ":") {
				return testSymbolObject(constant, actual[i])
			}

			if strings.HasPrefix(constant, "/") && strings.HasSuffix(constant, "/") {
				return testRegexpObject(constant[1:len(constant)-1], actual[i])
			}

			if strings.HasPrefix(constant, "class:") {
				return testClassObject(constant[6:], actual[i])
			}

			if strings.HasPrefix(constant, "module:") {
				return testModuleObject(constant[7:], actual[i])
			}

			return testStringObject(constant, actual[i])
		case []bytecode.Instructions:
			fn, ok := actual[i].Heap.(*object.Block)
			if !ok {
				return fmt.Errorf("constant %d - not a function: %T",
					i, actual[i])
			}
			err := testInstructions(constant, fn.Instructions)
			if err != nil {
				return fmt.Errorf("constant %d - testInstructions failed: %s",
					i, err)
			}
		}
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

func testStringObject(expected string, actual object.EmeraldValue) error {
	result, ok := actual.Heap.(*core.StringInstance)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}
	return nil
}

func testRegexpObject(expected string, actual object.EmeraldValue) error {
	result, ok := actual.Heap.(*core.RegexpInstance)
	if !ok {
		return fmt.Errorf("object is not Regexp. got=%T (%+v)", actual, actual)
	}
	if result.Source != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Source, expected)
	}
	return nil
}

func testSymbolObject(expected string, actual object.EmeraldValue) error {
	result, ok := actual.Heap.(*core.SymbolInstance)
	if !ok {
		return fmt.Errorf("object is not Symbol. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected[1:] {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}
	return nil
}

func testClassObject(expected string, actual object.EmeraldValue) error {
	class, ok := actual.Heap.(*object.Class)
	if !ok {
		return fmt.Errorf("object is not Class. got=%T (%+v)", actual, actual)
	}

	if class.Name != expected {
		return fmt.Errorf("class had wrong name want=%s, got=%s", expected, class.Name)
	}

	return nil
}

func testModuleObject(expected string, actual object.EmeraldValue) error {
	class, ok := actual.Heap.(*object.Module)
	if !ok {
		return fmt.Errorf("object is not Module. got=%T (%+v)", actual, actual)
	}

	if class.Name != expected {
		return fmt.Errorf("class had wrong name want=%s, got=%s", expected, class.Name)
	}

	return nil
}
