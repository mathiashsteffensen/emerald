package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileArrayLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []any{1, 2, 3},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpPushConstant, 5),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
