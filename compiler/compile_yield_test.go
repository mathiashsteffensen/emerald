package compiler

import "testing"

func TestCompileYield(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "yield 2, 5",
			expectedConstants: []any{2, 5},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpYield, 2),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
