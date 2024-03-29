package compiler

import "testing"

func TestCompileMethodLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `def method
				return 5 + 10
			end`,
			expectedConstants: []any{
				10,
				5,
				":method",
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpPushConstant, 1),
					Make(OpAdd),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			input: `def method
				5 + 10
			end`,
			expectedConstants: []any{
				10,
				5,
				":method",
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpPushConstant, 1),
					Make(OpAdd),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			input: `def method
				5
				10
			end`,
			expectedConstants: []any{
				5,
				10,
				":method",
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpPop),
					Make(OpPushConstant, 1),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			input: `def method; end`,
			expectedConstants: []any{
				":method",
				[]Instructions{
					Make(OpReturn),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			input: `
				def level=(new)
					@level = new
				end
			`,
			expectedConstants: []any{
				":@level",
				":level=",
				[]Instructions{
					Make(OpGetLocal, 0),
					Make(OpInstanceVarSet, 0),
					Make(OpReturnValue),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			input: `def method(name:); name end`,
			expectedConstants: []any{
				":method",
				[]Instructions{
					Make(OpGetLocal, 0),
					Make(OpReturn),
				},
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
