package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileSymbolLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             ":emerald",
			expectedConstants: []any{":emerald"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
