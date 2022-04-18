package compiler

import "testing"

func TestCompileAssignment(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			one = 1
			two = 2
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpSetGlobal, 1),
				Make(OpPop),
			},
		},
		{
			input: `
			one = 1
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpPop),
			},
		},
		{
			input: `
			one = 1
			two = one
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpSetGlobal, 1),
				Make(OpPop),
			},
		},
		{
			input: `
			two = one = 1
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpSetGlobal, 1),
				Make(OpPop),
			},
		},
		{
			input: `
			num = 55
			def method; num; end
			`,
			expectedConstants: []interface{}{
				55,
				":method",
				[]Instructions{
					Make(OpGetGlobal, 0),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			input: `
			def method
				num = 55
			end
			`,
			expectedConstants: []interface{}{
				55,
				":method",
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpSetLocal, 0),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpDefineMethod),
				Make(OpPop),
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
			expectedConstants: []interface{}{
				55,
				77,
				":method",
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpSetLocal, 0),
					Make(OpPop),
					Make(OpPushConstant, 1),
					Make(OpSetLocal, 1),
					Make(OpPop),
					Make(OpGetLocal, 1),
					Make(OpGetLocal, 0),
					Make(OpAdd),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			name:              "boolean and evaluating to true and assigning to variable",
			input:             "true && var = 15",
			expectedConstants: []any{15},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpJumpNotTruthy, 14),
				Make(OpPop),
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpJump, 15),
				Make(OpTrue),
				Make(OpPop),
			},
		},
		{
			name:              "boolean or evaluating to false and assigning to variable",
			input:             "false || var = 15",
			expectedConstants: []any{15},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpJumpNotTruthy, 9),
				Make(OpPop),
				Make(OpFalse),
				Make(OpJump, 15),
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpPop),
			},
		},
		{
			name:              "boolean and evaluating to true and assigning to variable",
			input:             "var = true; var &&= 15",
			expectedConstants: []any{15},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpSetGlobal, 0),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpJumpNotTruthy, 21),
				Make(OpPop),
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpJump, 24),
				Make(OpGetGlobal, 0),
				Make(OpPop),
			},
		},
		{
			name:              "boolean or evaluating to false and assigning to variable",
			input:             "var = false; var ||= 15",
			expectedConstants: []any{15},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpSetGlobal, 0),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpJumpNotTruthy, 18),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpJump, 24),
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpPop),
			},
		},
		{
			name:              "boolean or evaluating to false and assigning to an instance variable",
			input:             "@var = false; @var ||= 15",
			expectedConstants: []any{":@var", ":@var", ":@var", 15, ":@var"},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpInstanceVarSet, 0),
				Make(OpPop),
				Make(OpInstanceVarGet, 1),
				Make(OpJumpNotTruthy, 18),
				Make(OpPop),
				Make(OpInstanceVarGet, 2),
				Make(OpJump, 24),
				Make(OpPushConstant, 3),
				Make(OpInstanceVarSet, 4),
				Make(OpPop),
			},
		},
		{
			name:              "boolean and evaluating to true and assigning to an instance variable",
			input:             "@var = true; @var &&= 15",
			expectedConstants: []any{":@var", ":@var", 15, ":@var", ":@var"},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpInstanceVarSet, 0),
				Make(OpPop),
				Make(OpInstanceVarGet, 1),
				Make(OpJumpNotTruthy, 21),
				Make(OpPop),
				Make(OpPushConstant, 2),
				Make(OpInstanceVarSet, 3),
				Make(OpJump, 24),
				Make(OpInstanceVarGet, 4),
				Make(OpPop),
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
				[]Instructions{
					Make(OpInstanceVarGet, 0),
					Make(OpJumpNotTruthy, 13),
					Make(OpPop),
					Make(OpInstanceVarGet, 1),
					Make(OpJump, 19),
					Make(OpPushConstant, 2),
					Make(OpInstanceVarSet, 3),
					Make(OpReturnValue),
				},
				":get",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 4),
				Make(OpPushConstant, 5),
				Make(OpDefineMethod),
				Make(OpPop),
				Make(OpPushConstant, 6),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
