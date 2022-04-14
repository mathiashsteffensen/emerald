package compiler

import "testing"

func TestCompileHashLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "with no keys",
			input:             "{}",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				Make(OpHash, 0),
				Make(OpPop),
			},
		},
		{
			name:              "with constant keys and values",
			input:             "{1: 2, 3: 4, 5: 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
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
			input:             "{1: 2 + 3, 4: 5 * 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
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
