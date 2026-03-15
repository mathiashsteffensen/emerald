package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileBooleanLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "true",
			input:             "true",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "false",
			input:             "false",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
