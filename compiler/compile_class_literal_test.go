package compiler

import "testing"

func TestCompileClassLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:  "defining a new class",
			input: "class MyClass; end",
			expectedConstants: []any{
				":MyClass",
				"class:MyClass",
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGetOrSet, 0, 1),
				Make(OpOpenClass),
				Make(OpNull),
				Make(OpCloseClass),
				Make(OpPop),
			},
		},
		{
			name: "defining a namespaced class",
			input: "module MyMod; class MyClass; end; end\n" +
				"MyMod::MyClass",
			expectedConstants: []any{
				":MyMod",
				"module:MyMod",
				":MyClass",
				"class:MyClass",
				":MyMod",
				":MyClass",
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGetOrSet, 0, 1),
				Make(OpOpenClass),
				Make(OpConstantGetOrSet, 2, 3),
				Make(OpOpenClass),
				Make(OpNull),
				Make(OpCloseClass),
				Make(OpCloseClass),
				Make(OpPop),
				Make(OpConstantGet, 4),
				Make(OpScopedConstantGet, 5),
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
				":MyClass",
				"class:MyClass",
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
				Make(OpConstantGetOrSet, 0, 1),
				Make(OpOpenClass),
				Make(OpPushConstant, 3),
				Make(OpPushConstant, 4),
				Make(OpDefineMethod),
				Make(OpPop),
				Make(OpStaticTrue),
				Make(OpPushConstant, 6),
				Make(OpPushConstant, 7),
				Make(OpDefineMethod),
				Make(OpStaticFalse),
				Make(OpCloseClass),
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
				":MyClass",
				"class:MyClass",
				":MyClass",
				"class:MyClass",
				10,
				":my_method",
				[]Instructions{
					Make(OpPushConstant, 4),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGetOrSet, 0, 1),
				Make(OpOpenClass),
				Make(OpNull),
				Make(OpCloseClass),
				Make(OpPop),
				Make(OpConstantGetOrSet, 2, 3),
				Make(OpOpenClass),
				Make(OpPushConstant, 5),
				Make(OpPushConstant, 6),
				Make(OpDefineMethod),
				Make(OpCloseClass),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
