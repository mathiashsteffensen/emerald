package compiler

import (
	"testing"
)

func TestCompileSymbolLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             ":emerald",
			expectedConstants: []any{":emerald"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
