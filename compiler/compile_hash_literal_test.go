package compiler

import "testing"

func TestCompileHashLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "with no keys",
			input:             "{}",
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpHash, 0),
				Make(OpPop),
			},
		},
		{
			name:              "with index accessor",
			input:             "{}[:key]",
			expectedConstants: []any{":[]", ":key"},
			expectedInstructions: []Instructions{
				Make(OpHash, 0),
				Make(OpPushConstant, 0),
				Make(OpNull),
				Make(OpPushConstant, 1),
				Make(OpSend, 1),
				Make(OpPop),
			},
		},
		{
			name:              "with syntactic sugar for symbol keys",
			input:             "{key: 2}",
			expectedConstants: []any{":key", 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpHash, 2),
				Make(OpPop),
			},
		},
		{
			name:              "with constant keys and values",
			input:             "{1 => 2, 3 => 4, 5 => 6}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpPushConstant, 4),
				Make(OpPushConstant, 5),
				Make(OpHash, 6),
				Make(OpPop),
			},
		},
		{
			name:              "with expression values",
			input:             "{1 => 2 + 3, 4 => 5 * 6}",
			expectedConstants: []any{1, 3, 2, 4, 6, 5},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpAdd),
				Make(OpPushConstant, 3),
				Make(OpPushConstant, 4),
				Make(OpPushConstant, 5),
				Make(OpMul),
				Make(OpHash, 4),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
