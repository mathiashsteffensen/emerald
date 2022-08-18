package compiler

import "testing"

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
				[]Instructions{
					Make(OpPushConstant, 0), // The literal "24"
					Make(OpReturnValue),
				},
				":no_arg",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 1), // The symbol name of the method
				Make(OpPushConstant, 2), // The compiled block
				Make(OpDefineMethod),
				Make(OpPop),
				Make(OpPushConstant, 3), // The symbol name of the method
				Make(OpNull),            // Null block
				Make(OpSend, 0),
				Make(OpPop),
			},
		},
		{
			name:  "global var",
			input: "def method; $var = 5; $var; end",
			expectedConstants: []any{
				5,
				":method",
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpSetGlobal, 0),
					Make(OpPop),
					Make(OpGetGlobal, 0),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 1), // The symbol name of the method
				Make(OpPushConstant, 2), // The compiled block
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
