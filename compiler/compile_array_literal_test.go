package compiler

import "testing"

func TestCompileArrayLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpArray, 0),
				Make(OpPop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []any{1, 2, 3},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpArray, 3),
				Make(OpPop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []any{2, 1, 4, 3, 6, 5},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpAdd),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpSub),
				Make(OpPushConstant, 4),
				Make(OpPushConstant, 5),
				Make(OpMul),
				Make(OpArray, 3),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
