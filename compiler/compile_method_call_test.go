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
				Make(OpNull),
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
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpResetExecutionContext),
				Make(OpPop),
			},
		},
		{
			name:  "passing a block to a builtin method",
			input: "[0].map { |i| i + 2 }",
			expectedConstants: []any{
				0,
				":map",
				2,
				[]Instructions{
					Make(OpPushConstant, 2),
					Make(OpGetLocal, 0),
					Make(OpAdd),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpArray, 1),
				Make(OpSetExecutionContext),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 3),
				Make(OpSend, 0),
				Make(OpResetExecutionContext),
				Make(OpPop),
			},
		},
		{
			name: "calling a method with a receiver in block passed to a builtin method",
			input: `
			class Math; end
			[0].map { |i| Math.instance.add_two(i) }`,
			expectedConstants: []any{
				"class:Math",
				0,
				":map",
				":instance",
				":add_two",
				[]Instructions{
					Make(OpGetGlobal, 0),
					Make(OpSetExecutionContext),
					Make(OpPushConstant, 3),
					Make(OpNull),
					Make(OpSend, 0),
					Make(OpResetExecutionContext),
					Make(OpSetExecutionContext),
					Make(OpPushConstant, 4),
					Make(OpNull),
					Make(OpGetLocal, 0),
					Make(OpSend, 1),
					Make(OpResetExecutionContext),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpOpenClass),
				Make(OpNull),
				Make(OpCloseClass),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpArray, 1),
				Make(OpSetExecutionContext),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 5),
				Make(OpSend, 0),
				Make(OpResetExecutionContext),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
