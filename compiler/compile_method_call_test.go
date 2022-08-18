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
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 1),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpResetExecutionTarget),
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
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 1),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpResetExecutionTarget),
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
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 1),
				Make(OpCloseBlock, 3),
				Make(OpSend, 0),
				Make(OpResetExecutionTarget),
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
					Make(OpSetExecutionTarget),
					Make(OpPushConstant, 3),
					Make(OpNull),
					Make(OpSend, 0),
					Make(OpResetExecutionTarget),
					Make(OpSetExecutionTarget),
					Make(OpPushConstant, 4),
					Make(OpNull),
					Make(OpGetLocal, 0),
					Make(OpSend, 1),
					Make(OpResetExecutionTarget),
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
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 2),
				Make(OpCloseBlock, 5),
				Make(OpSend, 0),
				Make(OpResetExecutionTarget),
				Make(OpPop),
			},
		},
		{
			name: "closure test",
			input: `
				class MyClass
					[:one, :two, :three].each do |lvl|
						define_method(lvl) do |other_val|
							lvl
						end
					end
				end

				MyClass.new.one(:two)
			`,
			expectedConstants: []any{
				"class:MyClass",
				":one",
				":two",
				":three",
				":each",
				":define_method",
				[]Instructions{
					Make(OpGetFree, 0),
					Make(OpReturnValue),
				},
				[]Instructions{
					Make(OpPushConstant, 5),
					Make(OpGetLocal, 0),
					Make(OpCloseBlock, 6, 1),
					Make(OpGetLocal, 0),
					Make(OpSend, 1),
					Make(OpReturnValue),
				},
				":new",
				":one",
				":two",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpOpenClass),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpArray, 3),
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 4),
				Make(OpCloseBlock, 7, 0),
				Make(OpSend),
				Make(OpResetExecutionTarget),
				Make(OpCloseClass),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 8),
				Make(OpNull),
				Make(OpSend),
				Make(OpResetExecutionTarget),
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 9),
				Make(OpNull),
				Make(OpPushConstant, 10),
				Make(OpSend, 1),
				Make(OpResetExecutionTarget),
				Make(OpPop),
			},
		},
		{
			name: "setting an instance var on custom class with boolean or",
			input: `
				class MyClass
					class << self
						def instance; @instance ||= new; end
					end
				end

				MyClass.instance
				MyClass.instance
			`,
			expectedConstants: []any{
				"class:MyClass",
				":@instance",
				":@instance",
				":new",
				":@instance",
				":instance",
				[]Instructions{
					Make(OpInstanceVarGet, 1),
					Make(OpJumpNotTruthy, 13),
					Make(OpPop),
					Make(OpInstanceVarGet, 2),
					Make(OpJump, 22),
					Make(OpPushConstant, 3),
					Make(OpNull),
					Make(OpSend, 0),
					Make(OpInstanceVarSet, 4),
					Make(OpReturnValue),
				},
				":instance",
				":instance",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpOpenClass),
				Make(OpDefinitionStaticTrue),
				Make(OpPushConstant, 5),
				Make(OpPushConstant, 6),
				Make(OpDefineMethod),
				Make(OpDefinitionStaticFalse),
				Make(OpCloseClass),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 7),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpResetExecutionTarget),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 8),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpResetExecutionTarget),
				Make(OpPop),
			},
		},
		{
			name: "assignment method call",
			input: `
				class Logger
					class << self
						def level=(new)
							@level = new
						end
					end
				end

				Logger.level = :debug
			`,
			expectedConstants: []any{
				"class:Logger",
				":@level",
				":level=",
				[]Instructions{
					Make(OpGetLocal, 0),
					Make(OpInstanceVarSet, 1),
					Make(OpReturnValue),
				},
				":level=",
				":debug",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpOpenClass),
				Make(OpDefinitionStaticTrue),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpDefineMethod),
				Make(OpDefinitionStaticFalse),
				Make(OpCloseClass),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpSetExecutionTarget),
				Make(OpPushConstant, 4),
				Make(OpNull),
				Make(OpPushConstant, 5),
				Make(OpSend, 1),
				Make(OpResetExecutionTarget),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
