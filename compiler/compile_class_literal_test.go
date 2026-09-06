package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileClassLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:  "defining a new class",
			input: "class MyClass; end",
			expectedConstants: []any{
				":Object",
				":MyClass",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpOpenClass, 1),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
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
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpOpenModule, 0),
				bytecode.Make(bytecode.OpConstantGet, 1),
				bytecode.Make(bytecode.OpOpenClass, 2),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstantGet, 3),
				bytecode.Make(bytecode.OpScopedConstantGet, 4),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 2),
					bytecode.Make(bytecode.OpReturnValue),
				},
				15,
				":my_method",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 5),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpOpenClass, 1),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpStaticTrue),
				bytecode.Make(bytecode.OpPushConstant, 6),
				bytecode.Make(bytecode.OpPushConstant, 7),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpStaticFalse),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 4),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstantGet, 0),
				bytecode.Make(bytecode.OpOpenClass, 1),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstantGet, 2),
				bytecode.Make(bytecode.OpOpenClass, 3),
				bytecode.Make(bytecode.OpPushConstant, 5),
				bytecode.Make(bytecode.OpPushConstant, 6),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
