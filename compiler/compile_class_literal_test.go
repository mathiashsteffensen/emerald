package compiler

import "testing"

func TestCompileClassLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:  "defining a new class",
			input: "class MyClass; end",
			expectedConstants: []any{
				"class:MyClass",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpOpenClass),
				Make(OpNull),
				Make(OpCloseClass),
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
				"class:MyClass",
				10,
				":my_method",
				[]Instructions{
					Make(OpPushConstant, 1),
					Make(OpReturnValue),
				},
				15,
				":my_method",
				[]Instructions{
					Make(OpPushConstant, 4),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpOpenClass),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpDefineMethod),
				Make(OpPop),
				Make(OpDefinitionStaticTrue),
				Make(OpPushConstant, 5),
				Make(OpPushConstant, 6),
				Make(OpDefineMethod),
				Make(OpDefinitionStaticFalse),
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
				"class:MyClass",
				10,
				":my_method",
				[]Instructions{
					Make(OpPushConstant, 1),
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
				Make(OpGetGlobal, 0),
				Make(OpOpenClass),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpDefineMethod),
				Make(OpCloseClass),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
