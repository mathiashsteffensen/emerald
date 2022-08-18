package compiler

import "testing"

func TestCompileStringLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"emerald"`,
			expectedConstants: []any{"emerald"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
