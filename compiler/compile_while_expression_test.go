package compiler

import (
	"emerald/bytecode"
	"testing"
)

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
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 3, 0, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpJumpNotTruthy, 46),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpGetGlobal, 1),
				bytecode.Make(bytecode.OpSend, 4, 1, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpJump, 16),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
				items = [1, 4, 9]
			
				puts("Hello") while items.pop
			`,
			expectedConstants: []any{1, 4, 9, ":pop", ":puts", "Hello"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 3, 0, 0),
				bytecode.Make(bytecode.OpJumpNotTruthy, 43),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPushConstant, 5),
				bytecode.Make(bytecode.OpSend, 4, 1, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpJump, 16),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
