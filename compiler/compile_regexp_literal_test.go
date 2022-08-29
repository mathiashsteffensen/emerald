package compiler

import "testing"

func TestCompileRegexpLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "/abc/.match(\"abc\")",
			expectedConstants: []any{
				"/abc/",
				":match",
				"abc",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpNull),
				Make(OpPushConstant, 2),
				Make(OpSend, 1),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
