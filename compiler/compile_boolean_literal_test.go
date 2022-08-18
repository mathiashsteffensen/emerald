package compiler

import "testing"

func TestCompileBooleanLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "true",
			input:             "true",
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpPop),
			},
		},
		{
			name:              "false",
			input:             "false",
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
