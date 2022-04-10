package compiler

import "testing"

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "simple addition",
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpAdd),
				Make(OpPop),
			},
		},
		{
			name:              "stack cleanup",
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpPop),
			},
		},
		{
			name:              "subtracting",
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpSub),
				Make(OpPop),
			},
		},
		{
			name:              "multiplying",
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpMul),
				Make(OpPop),
			},
		},
		{
			name:              "dividing",
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpDiv),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
