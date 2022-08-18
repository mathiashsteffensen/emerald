package compiler

import "testing"

func TestCompilePrefixExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "negating integer",
			input:             "-1",
			expectedConstants: []any{1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpMinus),
				Make(OpPop),
			},
		},
		{
			name:              "negating boolean",
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []Instructions{
				Make(OpTrue),
				Make(OpBang),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
