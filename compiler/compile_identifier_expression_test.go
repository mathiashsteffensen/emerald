package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileIdentifierExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			def no_arg
				24
			end
			no_arg
			`,
			expectedConstants: []any{
				24,
				":no_arg",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 0), // The literal "24"
					bytecode.Make(bytecode.OpReturnValue),
				},
				":no_arg",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 1), // The symbol name of the method
				bytecode.Make(bytecode.OpPushConstant, 2), // The compiled block
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpSelf), // Receiver, implicit self
				bytecode.Make(bytecode.OpNull), // Null block
				bytecode.Make(bytecode.OpSend, 3, 0, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:  "global var",
			input: "def method; $var = 5; $var; end",
			expectedConstants: []any{
				5,
				":method",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 0),
					bytecode.Make(bytecode.OpSetGlobal, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetGlobal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 1), // The symbol name of the method
				bytecode.Make(bytecode.OpPushConstant, 2), // The compiled block
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
