package compiler

import (
	"emerald/bytecode"
	"testing"
)

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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 0),
					bytecode.Make(bytecode.OpPushConstant, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 0),
					bytecode.Make(bytecode.OpPushConstant, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpPushConstant, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `def method; end`,
			expectedConstants: []any{
				":method",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
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
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpInstanceVarSet, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `def method(name:); name end`,
			expectedConstants: []any{
				":method",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturn),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
