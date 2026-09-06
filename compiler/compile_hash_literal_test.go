package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileHashLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "with no keys",
			input:             "{}",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpHash, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "with index accessor",
			input:             "{}[:key]",
			expectedConstants: []any{":[]", ":key"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpHash, 0),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpSend, 0, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "with syntactic sugar for symbol keys",
			input:             "{key: 2}",
			expectedConstants: []any{":key", 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpHash, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "with constant keys and values",
			input:             "{1 => 2, 3 => 4, 5 => 6}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpPushConstant, 5),
				bytecode.Make(bytecode.OpHash, 6),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "with expression values",
			input:             "{1 => 2 + 3, 4 => 5 * 6}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpPushConstant, 5),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpHash, 4),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
