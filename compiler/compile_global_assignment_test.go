package compiler

import "testing"

func TestGlobalAssignments(t *testing.T) {
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
				Make(OpGetGlobal, 0),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpSetGlobal, 1),
				Make(OpGetGlobal, 1),
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
				Make(OpGetGlobal, 0),
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
				Make(OpGetGlobal, 0),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpSetGlobal, 1),
				Make(OpGetGlobal, 1),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
