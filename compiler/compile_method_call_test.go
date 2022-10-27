package compiler

import (
	"testing"
)

func TestCompileMethodCall(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "static receiver",
			input: `
			class MyClass; end

			MyClass.new
			`,
			expectedConstants: []any{
				":Object",
				":MyClass",
				":MyClass",
				":new",
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpOpenClass, 1),
				Make(OpNull),
				Make(OpUnwrapContext),
				Make(OpPop),
				Make(OpConstantGet, 2),
				Make(OpPushConstant, 3),
				Make(OpNull),
				Make(OpSend, 0),
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
				Make(OpPushConstant, 1),
				Make(OpNull),
				Make(OpSend, 0),
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
				Make(OpPushConstant, 1),
				Make(OpCloseBlock, 3),
				Make(OpSend, 0),
				Make(OpPop),
			},
		},
		{
			name: "calling a method with a receiver in block",
			input: `
			class Math; end
			[0].map { |i| Math.instance.add_two(i) }`,
			expectedConstants: []any{
				":Object",
				":Math",
				0,
				":map",
				":Math",
				":instance",
				":add_two",
				[]Instructions{
					Make(OpConstantGet, 4),
					Make(OpPushConstant, 5),
					Make(OpNull),
					Make(OpSend, 0),
					Make(OpPushConstant, 6),
					Make(OpNull),
					Make(OpGetLocal, 0),
					Make(OpSend, 1),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpOpenClass, 1),
				Make(OpNull),
				Make(OpUnwrapContext),
				Make(OpPop),
				Make(OpPushConstant, 2),
				Make(OpArray, 1),
				Make(OpPushConstant, 3),
				Make(OpCloseBlock, 7),
				Make(OpSend, 0),
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
				":Object",
				":MyClass",
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
					Make(OpSelf),
					Make(OpPushConstant, 6),
					Make(OpGetLocal, 0),
					Make(OpCloseBlock, 7, 1),
					Make(OpGetLocal, 0),
					Make(OpSend, 1),
					Make(OpReturnValue),
				},
				":MyClass",
				":new",
				":one",
				":two",
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpOpenClass, 1),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpPushConstant, 4),
				Make(OpArray, 3),
				Make(OpPushConstant, 5),
				Make(OpCloseBlock, 8, 0),
				Make(OpSend),
				Make(OpUnwrapContext),
				Make(OpPop),
				Make(OpConstantGet, 9),
				Make(OpPushConstant, 10),
				Make(OpNull),
				Make(OpSend),
				Make(OpPushConstant, 11),
				Make(OpNull),
				Make(OpPushConstant, 12),
				Make(OpSend, 1),
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
				":Object",
				":MyClass",
				":@instance",
				":@instance",
				":new",
				":@instance",
				":instance",
				[]Instructions{
					Make(OpInstanceVarGet, 2),
					Make(OpJumpNotTruthy, 13),
					Make(OpPop),
					Make(OpInstanceVarGet, 3),
					Make(OpJump, 23),
					Make(OpSelf),
					Make(OpPushConstant, 4),
					Make(OpNull),
					Make(OpSend, 0),
					Make(OpInstanceVarSet, 5),
					Make(OpReturnValue),
				},
				":MyClass",
				":instance",
				":MyClass",
				":instance",
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpOpenClass, 1),
				Make(OpStaticTrue),
				Make(OpPushConstant, 6),
				Make(OpPushConstant, 7),
				Make(OpDefineMethod),
				Make(OpStaticFalse),
				Make(OpUnwrapContext),
				Make(OpPop),
				Make(OpConstantGet, 8),
				Make(OpPushConstant, 9),
				Make(OpNull),
				Make(OpSend, 0),
				Make(OpPop),
				Make(OpConstantGet, 10),
				Make(OpPushConstant, 11),
				Make(OpNull),
				Make(OpSend, 0),
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
				":Object",
				":Logger",
				":@level",
				":level=",
				[]Instructions{
					Make(OpGetLocal, 0),
					Make(OpInstanceVarSet, 2),
					Make(OpReturnValue),
				},
				":Logger",
				":level=",
				":debug",
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpOpenClass, 1),
				Make(OpStaticTrue),
				Make(OpPushConstant, 3),
				Make(OpPushConstant, 4),
				Make(OpDefineMethod),
				Make(OpStaticFalse),
				Make(OpUnwrapContext),
				Make(OpPop),
				Make(OpConstantGet, 5),
				Make(OpPushConstant, 6),
				Make(OpNull),
				Make(OpPushConstant, 7),
				Make(OpSend, 1),
				Make(OpPop),
			},
		},
		{
			name:  "keyword arguments",
			input: "puts one: :two, three: :four",
			expectedConstants: []any{
				":puts",
				":one",
				":two",
				":three",
				":four",
			},
			expectedInstructions: []Instructions{
				Make(OpSelf),
				Make(OpPushConstant, 0),
				Make(OpNull),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpPushConstant, 4),
				Make(OpHash, 4),
				Make(OpSend, 2, 1),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
