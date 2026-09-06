package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileAssignment(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				one = 1
				two = 2
			`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
				One = 1
			`,
			expectedConstants: []any{1, ":One"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpConstantSet, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			one = 1
			`,
			expectedConstants: []any{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			one = 1
			two = one
			`,
			expectedConstants: []any{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			two = one = 1
			`,
			expectedConstants: []any{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			num = 55
			def method; num; end
			`,
			expectedConstants: []any{
				55,
				":method",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetGlobal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			def method
				num = 55
			end
			`,
			expectedConstants: []any{
				55,
				":method",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			def method
				a = 55
				b = 77
				a + b
			end
			`,
			expectedConstants: []any{
				55,
				77,
				":method",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpPushConstant, 1),
					bytecode.Make(bytecode.OpSetLocal, 1),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean and evaluating to true and assigning to variable 1",
			input:             "true && var = 15",
			expectedConstants: []any{15},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJumpNotTruthy, 14),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpJump, 15),
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean or evaluating to false and assigning to variable 1",
			input:             "false || var = 15",
			expectedConstants: []any{15},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpJumpNotTruthy, 9),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpJump, 15),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean and evaluating to true and assigning to variable 2",
			input:             "var = true; var &&= 15",
			expectedConstants: []any{15},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpJumpNotTruthy, 21),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpJump, 24),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean or evaluating to false and assigning to variable 2",
			input:             "var = false; var ||= 15",
			expectedConstants: []any{15},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpJumpNotTruthy, 18),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpJump, 24),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean or evaluating to false and assigning to an instance variable",
			input:             "@var = false; @var ||= 15",
			expectedConstants: []any{":@var", ":@var", ":@var", 15, ":@var"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpInstanceVarSet, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpInstanceVarGet, 1),
				bytecode.Make(bytecode.OpJumpNotTruthy, 18),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpInstanceVarGet, 2),
				bytecode.Make(bytecode.OpJump, 24),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpInstanceVarSet, 4),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "boolean and evaluating to true and assigning to an instance variable",
			input:             "@var = true; @var &&= 15",
			expectedConstants: []any{":@var", ":@var", 15, ":@var", ":@var"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpInstanceVarSet, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpInstanceVarGet, 1),
				bytecode.Make(bytecode.OpJumpNotTruthy, 21),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpInstanceVarSet, 3),
				bytecode.Make(bytecode.OpJump, 24),
				bytecode.Make(bytecode.OpInstanceVarGet, 4),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:  "returning boolean or from method",
			input: "def get; @var ||= 5; end; get",
			expectedConstants: []any{
				":@var",
				":@var",
				5,
				":@var",
				":get",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpInstanceVarGet, 0),
					bytecode.Make(bytecode.OpJumpNotTruthy, 13),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpInstanceVarGet, 1),
					bytecode.Make(bytecode.OpJump, 19),
					bytecode.Make(bytecode.OpPushConstant, 2),
					bytecode.Make(bytecode.OpInstanceVarSet, 3),
					bytecode.Make(bytecode.OpReturnValue),
				},
				":get",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpPushConstant, 5),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 6, 0, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "assigning to global var",
			input:             "$foo = 5",
			expectedConstants: []any{5},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
