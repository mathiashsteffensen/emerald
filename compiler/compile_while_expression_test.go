package compiler

import "testing"

func TestCompileWhileExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				items = [1, 4, 9]
				while item = items.pop
					puts(item)
				end
			`,
			expectedConstants: []any{1, 4, 9, ":pop", ":puts"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpArray, 3),
				Make(OpSetGlobal, 0),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpPushConstant, 3),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpSetGlobal, 1),
				Make(OpJumpNotTruthy, 46),
				Make(OpPop),
				Make(OpSelf),
				Make(OpPushConstant, 4),
				Make(OpNull),
				Make(OpGetGlobal, 1),
				Make(OpSend, 1),
				Make(OpPop),
				Make(OpJump, 16),
				Make(OpNull),
				Make(OpPop),
			},
		},
		{
			input: `
				items = [1, 4, 9]
			
				puts("Hello") while items.pop
			`,
			expectedConstants: []any{1, 4, 9, ":pop", ":puts", "Hello"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpArray, 3),
				Make(OpSetGlobal, 0),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpPushConstant, 3),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpJumpNotTruthy, 43),
				Make(OpPop),
				Make(OpSelf),
				Make(OpPushConstant, 4),
				Make(OpNull),
				Make(OpPushConstant, 5),
				Make(OpSend, 1),
				Make(OpPop),
				Make(OpJump, 16),
				Make(OpNull),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
