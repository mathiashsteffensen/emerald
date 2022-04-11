package compiler

import "testing"

func TestCompileStringLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"emerald"`,
			expectedConstants: []interface{}{"emerald"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
