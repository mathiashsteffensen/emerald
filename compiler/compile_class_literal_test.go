package compiler

import "testing"

func TestCompileClassLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:  "defining a new class",
			input: "class MyClass; end",
			expectedConstants: []any{
				":Object",
				":MyClass",
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpOpenClass, 1),
				Make(OpNull),
				Make(OpUnwrapContext),
				Make(OpPop),
			},
		},
		{
			name: "defining a namespaced class",
			input: `
				module MyMod
					class MyClass; end
				end
				
				MyMod::MyClass
			`,
			expectedConstants: []any{
				":MyMod",
				":Object",
				":MyClass",
				":MyMod",
				":MyClass",
			},
			expectedInstructions: []Instructions{
				Make(OpOpenModule, 0),
				Make(OpConstantGet, 1),
				Make(OpOpenClass, 2),
				Make(OpNull),
				Make(OpUnwrapContext),
				Make(OpUnwrapContext),
				Make(OpPop),
				Make(OpConstantGet, 3),
				Make(OpScopedConstantGet, 4),
				Make(OpPop),
			},
		},
		{
			name: "defining a new class with a body",
			input: `
			class MyClass
				def my_method
					10
				end

				class << self
					def my_method
						15
					end
				end
			end
			`,
			expectedConstants: []any{
				":Object",
				":MyClass",
				10,
				":my_method",
				[]Instructions{
					Make(OpPushConstant, 2),
					Make(OpReturnValue),
				},
				15,
				":my_method",
				[]Instructions{
					Make(OpPushConstant, 5),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpOpenClass, 1),
				Make(OpPushConstant, 3),
				Make(OpPushConstant, 4),
				Make(OpDefineMethod),
				Make(OpPop),
				Make(OpStaticTrue),
				Make(OpPushConstant, 6),
				Make(OpPushConstant, 7),
				Make(OpDefineMethod),
				Make(OpStaticFalse),
				Make(OpUnwrapContext),
				Make(OpPop),
			},
		},
		{
			name: "overwriting an existing class",
			input: `
			class MyClass; end
			class MyClass
				def my_method
					10
				end
			end
			`,
			expectedConstants: []any{
				":Object",
				":MyClass",
				":Object",
				":MyClass",
				10,
				":my_method",
				[]Instructions{
					Make(OpPushConstant, 4),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpOpenClass, 1),
				Make(OpNull),
				Make(OpUnwrapContext),
				Make(OpPop),
				Make(OpConstantGet, 2),
				Make(OpOpenClass, 3),
				Make(OpPushConstant, 5),
				Make(OpPushConstant, 6),
				Make(OpDefineMethod),
				Make(OpUnwrapContext),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
