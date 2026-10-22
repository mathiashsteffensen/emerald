package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileBreakStatement(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				while true
					break 1
				end
			`,
			expectedConstants: []any{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJumpNotTruthy, 14),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpJump, 15),
				bytecode.Make(bytecode.OpJump, 0),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
