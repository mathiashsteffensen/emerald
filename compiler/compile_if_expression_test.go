package compiler

import "testing"

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if true
				10 
			end
			3333
			`,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []Instructions{
				// 0000
				Make(OpTrue),
				// 0001
				Make(OpJumpNotTruthy, 10),
				// 0004
				Make(OpPushConstant, 0),
				// 0007
				Make(OpJump, 11),
				// 0010
				Make(OpNull),
				// 0011
				Make(OpPop),
				// 0012
				Make(OpPushConstant, 1),
				// 0015
				Make(OpPop),
			},
		},
		{
			input: `
			if true
				10
			else
				20
			end
			3333
			`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []Instructions{
				// 0000
				Make(OpTrue),
				// 0001
				Make(OpJumpNotTruthy, 10),
				// 0004
				Make(OpPushConstant, 0),
				// 0007
				Make(OpJump, 13),
				// 0010
				Make(OpPushConstant, 1),
				// 0013
				Make(OpPop),
				// 0014
				Make(OpPushConstant, 2),
				// 0017
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
