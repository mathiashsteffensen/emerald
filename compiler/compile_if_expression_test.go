package compiler

import "testing"

func TestCompileIfExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if true
				10 
			end
			3333
			`,
			expectedConstants: []any{10, 3333},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpJumpNotTruthy, 11),
				Make(OpPop),
				Make(OpPushConstant, 0),
				Make(OpJump, 12),
				Make(OpNull),
				Make(OpPop),
				Make(OpPushConstant, 1),
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
			expectedConstants: []any{10, 20, 3333},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpJumpNotTruthy, 11),
				Make(OpPop),
				Make(OpPushConstant, 0),
				Make(OpJump, 14),
				Make(OpPushConstant, 1),
				Make(OpPop),
				Make(OpPushConstant, 2),
				Make(OpPop),
			},
		},
		{
			name:              "negate if expression resolving to nil",
			input:             "!(if false; 5; end)",
			expectedConstants: []any{5},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpJumpNotTruthy, 11),
				Make(OpPop),
				Make(OpPushConstant, 0),
				Make(OpJump, 12),
				Make(OpNull),
				Make(OpBang),
				Make(OpPop),
			},
		},
		{
			name: "with elsif",
			input: `
				if true
					true
				elsif true
					true
				elsif false
					true
				end
			`,
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpJumpNotTruthy, 9),
				Make(OpPop),
				Make(OpTrue),
				Make(OpJump, 28),
				Make(OpTrue),
				Make(OpJumpNotTruthy, 18),
				Make(OpPop),
				Make(OpTrue),
				Make(OpJump, 28),
				Make(OpFalse),
				Make(OpJumpNotTruthy, 27),
				Make(OpPop),
				Make(OpTrue),
				Make(OpJump, 28),
				Make(OpNull),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
