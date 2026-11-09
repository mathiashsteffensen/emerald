package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileUnicodeIdentifiers(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "uppercase non-ASCII first letter assigns a constant",
			input:             "Ärger = 1",
			expectedConstants: []any{1, ":Ärger"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpConstantSet, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "titlecase first letter assigns a constant",
			input:             "ǅx = 1",
			expectedConstants: []any{1, ":ǅx"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpConstantSet, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			// The first byte of é is 0xC3, which is also the code point of uppercase Ã
			name:              "lowercase non-ASCII first letter assigns a local",
			input:             "éclair = 1; éclair",
			expectedConstants: []any{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:              "uppercase non-ASCII first letter reads a constant",
			input:             "Ärger",
			expectedConstants: []any{":Ärger"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
