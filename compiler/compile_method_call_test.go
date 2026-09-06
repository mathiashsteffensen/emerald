package compiler

import (
	"emerald/bytecode"
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
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpOpenClass, 1),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstantGet, 2),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 3, 0, 0),
				bytecode.Make(bytecode.OpPop),
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
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 1, 0, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:  "passing a block to a builtin method",
			input: "[0].map { |i| i + 2 }",
			expectedConstants: []any{
				0,
				":map",
				2,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpPushConstant, 2),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpArray, 1),
				bytecode.Make(bytecode.OpCloseBlock, 3),
				bytecode.Make(bytecode.OpSend, 1, 0, 0),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstantGet, 4),
					bytecode.Make(bytecode.OpNull),
					bytecode.Make(bytecode.OpSend, 5, 0, 0),
					bytecode.Make(bytecode.OpNull),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpSend, 6, 1, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpOpenClass, 1),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpArray, 1),
				bytecode.Make(bytecode.OpCloseBlock, 7),
				bytecode.Make(bytecode.OpSend, 3, 0, 0),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetFree, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpSelf),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpCloseBlock, 7, 1),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpSend, 6, 1, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
				":MyClass",
				":new",
				":one",
				":two",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpOpenClass, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpCloseBlock, 8, 0),
				bytecode.Make(bytecode.OpSend, 5, 0, 0),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstantGet, 9),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 10, 0, 0),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPushConstant, 12),
				bytecode.Make(bytecode.OpSend, 11, 1, 0),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpInstanceVarGet, 2),
					bytecode.Make(bytecode.OpJumpNotTruthy, 13),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpInstanceVarGet, 3),
					bytecode.Make(bytecode.OpJump, 23),
					bytecode.Make(bytecode.OpSelf),
					bytecode.Make(bytecode.OpNull),
					bytecode.Make(bytecode.OpSend, 4, 0, 0),
					bytecode.Make(bytecode.OpInstanceVarSet, 5),
					bytecode.Make(bytecode.OpReturnValue),
				},
				":MyClass",
				":instance",
				":MyClass",
				":instance",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpOpenClass, 1),
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpStaticTrue),
				bytecode.Make(bytecode.OpPushConstant, 6),
				bytecode.Make(bytecode.OpPushConstant, 7),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpStaticFalse),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstantGet, 8),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 9, 0, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstantGet, 10),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpSend, 11, 0, 0),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpInstanceVarSet, 2),
					bytecode.Make(bytecode.OpReturnValue),
				},
				":Logger",
				":level=",
				":debug",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpOpenClass, 1),
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpStaticTrue),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpStaticFalse),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstantGet, 5),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPushConstant, 7),
				bytecode.Make(bytecode.OpSendAssign, 6, 1, 0),
				bytecode.Make(bytecode.OpPop),
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
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpHash, 4),
				bytecode.Make(bytecode.OpSend, 0, 2, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
