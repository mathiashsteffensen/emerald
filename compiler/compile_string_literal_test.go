package compiler

import "testing"

func TestCompileStringLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"emerald"`,
			expectedConstants: []any{"emerald"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPop),
			},
		},
		{
			input:             `placeholder = "template"; "This is a #{placeholder.to_s}"`,
			expectedConstants: []any{"template", "This is a", ":to_s"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpGetGlobal, 0),
				Make(OpPushConstant, 2),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpStringJoin, 2),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
