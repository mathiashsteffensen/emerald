package compiler

import "testing"

func TestBooleanLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "true",
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpPop),
			},
		},
		{
			name:              "false",
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []Instructions{
				Make(OpFalse),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
