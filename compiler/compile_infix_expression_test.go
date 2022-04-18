package compiler

import "testing"

func TestCompileInfixExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "addition",
			input:             "1 + 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpAdd),
				Make(OpPop),
			},
		},
		{
			name:              "subtracting",
			input:             "1 - 2",
			expectedConstants: []any{2, 1},
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
			expectedConstants: []any{2, 1},
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
			expectedConstants: []any{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpDiv),
				Make(OpPop),
			},
		},
		{
			name:              "stack cleanup",
			input:             "1; 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpPop),
			},
		},
		{
			name:              "adding strings",
			input:             `"eme" + "rald"`,
			expectedConstants: []any{"rald", "eme"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpAdd),
				Make(OpPop),
			},
		},
		{
			name:              "greater than",
			input:             "1 > 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpGreaterThan),
				Make(OpPop),
			},
		},
		{
			name:              "greater than or eq",
			input:             "1 >= 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpGreaterThanOrEq),
				Make(OpPop),
			},
		},
		{
			name:              "less than",
			input:             "1 < 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpLessThan),
				Make(OpPop),
			},
		},
		{
			name:              "less than or eq",
			input:             "1 <= 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpLessThanOrEq),
				Make(OpPop),
			},
		},
		{
			name:              "integers equals",
			input:             "1 == 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpEqual),
				Make(OpPop),
			},
		},
		{
			name:              "integers not equals",
			input:             "1 != 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpNotEqual),
				Make(OpPop),
			},
		},
		{
			name:              "boolean equals",
			input:             "true == false",
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpTrue),
				Make(OpEqual),
				Make(OpPop),
			},
		},
		{
			name:              "boolean not equals",
			input:             "true != false",
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpTrue),
				Make(OpNotEqual),
				Make(OpPop),
			},
		},
		{
			name:              "boolean and",
			input:             "true && false",
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpJumpNotTruthy, 7),
				Make(OpPop),
				Make(OpFalse),
				Make(OpPop),
			},
		},
		{
			name:              "boolean or",
			input:             "1 + 2 || false",
			expectedConstants: []any{2, 1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpAdd),
				Make(OpJumpTruthy, 13),
				Make(OpPop),
				Make(OpFalse),
				Make(OpPop),
			},
		},
		{
			name:              "boolean and evaluating to false",
			input:             "false && 15",
			expectedConstants: []any{15},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpJumpNotTruthy, 9),
				Make(OpPop),
				Make(OpPushConstant, 0),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
