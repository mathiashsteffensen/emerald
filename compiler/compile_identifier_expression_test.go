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
			expectedConstants: []interface{}{
				24,
				[]Instructions{
					Make(OpPushConstant, 0), // The literal "24"
					Make(OpReturnValue),
				},
				":no_arg",
				":no_arg",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 1), // The compiled block
				Make(OpPushConstant, 2), // The symbol name of the method
				Make(OpDefineMethod),
				Make(OpPop),
				Make(OpPushConstant, 3), // The symbol name of the method
				Make(OpSend, 0),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
