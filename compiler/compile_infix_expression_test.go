package compiler

import "testing"

func TestCompileInfixExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "addition",
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
			name:              "adding strings",
			input:             `"eme" + "rald"`,
			expectedConstants: []interface{}{"eme", "rald"},
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
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpGreaterThan),
				Make(OpPop),
			},
		},
		{
			name:              "less than",
			input:             "1 < 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpLessThan),
				Make(OpPop),
			},
		},
		{
			name:              "integers equals",
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
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
			expectedConstants: []interface{}{1, 2},
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
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpFalse),
				Make(OpEqual),
				Make(OpPop),
			},
		},
		{
			name:              "boolean not equals",
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpFalse),
				Make(OpNotEqual),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
