package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompilePrefixExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "negating integer",
			input:             "-1",
			expectedConstants: []any{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpMinus),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "negating boolean",
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpBang),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
