package compiler

import "testing"

func TestCompileCallExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			def one_arg(a); a; end
			one_arg(24)
			`,
			expectedConstants: []any{
				":one_arg",
				[]Instructions{
					Make(OpGetLocal, 0),
					Make(OpReturnValue),
				},
				":one_arg",
				24,
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpDefineMethod),
				Make(OpPop),
				Make(OpSelf),
				Make(OpPushConstant, 2),
				Make(OpNull),
				Make(OpPushConstant, 3),
				Make(OpSend, 1),
				Make(OpPop),
			},
		},
		{
			input: `
			def many_arg(a, b, c); a; b; c; end
			many_arg(24, 25, 26)
			`,
			expectedConstants: []any{
				":many_arg",
				[]Instructions{
					Make(OpGetLocal, 0),
					Make(OpPop),
					Make(OpGetLocal, 1),
					Make(OpPop),
					Make(OpGetLocal, 2),
					Make(OpReturnValue),
				},
				":many_arg",
				24,
				25,
				26,
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpDefineMethod),
				Make(OpPop),
				Make(OpSelf),
				Make(OpPushConstant, 2),
				Make(OpNull),
				Make(OpPushConstant, 3),
				Make(OpPushConstant, 4),
				Make(OpPushConstant, 5),
				Make(OpSend, 3),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
