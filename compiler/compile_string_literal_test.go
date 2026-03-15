package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileStringLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"emerald"`,
			expectedConstants: []any{"emerald"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             `placeholder = "template"; "This is a #{placeholder.to_s}"`,
			expectedConstants: []any{"template", "This is a", ":to_s"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 0),
				bytecode.Make(bytecode.OpStringJoin, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
