package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileYield(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "yield 2, 5",
			expectedConstants: []any{2, 5},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpYield, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
