package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileInfixExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "addition",
			input:             "1.0 + 2",
			expectedConstants: []any{1.0, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "subtracting",
			input:             "1 - 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "multiplying",
			input:             "1 * 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "dividing",
			input:             "2 / 1",
			expectedConstants: []any{2, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpDiv),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "stack cleanup",
			input:             "1; 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "adding strings",
			input:             `"eme" + "rald"`,
			expectedConstants: []any{"eme", "rald"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "greater than",
			input:             "1 > 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "greater than or eq",
			input:             "1 >= 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpGreaterThanOrEq),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "less than",
			input:             "1 < 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpLessThan),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "less than or eq",
			input:             "1 <= 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpLessThanOrEq),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "integers equals",
			input:             "1 == 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "integers not equals",
			input:             "1 != 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean equals",
			input:             "true == false",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean not equals",
			input:             "true != false",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean and",
			input:             "true && false",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpDupN, 1),
				bytecode.Make(bytecode.OpJumpNotTruthy, 9),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean or",
			input:             "1 + 2 || false",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpDupN, 1),
				bytecode.Make(bytecode.OpJumpTruthy, 15),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean and evaluating to false",
			input:             "false && 15",
			expectedConstants: []any{15},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpDupN, 1),
				bytecode.Make(bytecode.OpJumpNotTruthy, 11),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "spaceship",
			input:             "1 <=> 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpSpaceship),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "matching",
			input:             `/a/ =~ "a"`,
			expectedConstants: []any{"/a/", "a"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpMatch),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "case equality",
			input:             `true === true`,
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpCaseEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "binary shift left",
			input:             "[] << 2",
			expectedConstants: []any{2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpArray),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpBinShiftLeft),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
