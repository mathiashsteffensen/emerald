package compiler

import (
	"emerald/bytecode"
	"testing"
)

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
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJumpNotTruthy, 11),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpJump, 12),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPop),
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
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJumpNotTruthy, 11),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpJump, 14),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "negate if expression resolving to nil",
			input:             "!(if false; 5; end)",
			expectedConstants: []any{5},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpJumpNotTruthy, 11),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpJump, 12),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpBang),
				bytecode.Make(bytecode.OpPop),
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
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJumpNotTruthy, 9),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJump, 28),
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJumpNotTruthy, 18),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJump, 28),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpJumpNotTruthy, 27),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJump, 28),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
