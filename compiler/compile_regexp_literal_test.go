package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileRegexpLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "/abc/.match(\"abc\")",
			expectedConstants: []any{
				"/abc/",
				":match",
				"abc",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpSend, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
