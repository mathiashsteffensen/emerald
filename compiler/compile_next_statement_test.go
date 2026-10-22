package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileNextStatement(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				while true
					next
				end
			`,
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJumpNotTruthy, 13),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpJump, 0),
				bytecode.Make(bytecode.OpJump, 0),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
