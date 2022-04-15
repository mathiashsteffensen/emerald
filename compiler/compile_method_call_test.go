package compiler

import "testing"

func TestCompileMethodCall(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "static receiver",
			input: `
			class MyClass; end

			MyClass.new
			`,
			expectedConstants: []any{
				"class:MyClass",
				":new",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpOpenClass),
				Make(OpNull),
				Make(OpCloseClass),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpSetExecutionContext),
				Make(OpPushConstant, 1),
				Make(OpSend, 0),
				Make(OpResetExecutionContext),
				Make(OpPop),
			},
		},
		{
			name: "instance receiver",
			input: `
			"string".to_sym
			`,
			expectedConstants: []any{
				"string",
				":to_sym",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetExecutionContext),
				Make(OpPushConstant, 1),
				Make(OpSend, 0),
				Make(OpResetExecutionContext),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
