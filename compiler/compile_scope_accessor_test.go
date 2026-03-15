package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileScopeAccessor(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "MyMod::MyClass",
			expectedConstants: []any{
				":MyMod",
				":MyClass",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpScopedConstantGet, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
