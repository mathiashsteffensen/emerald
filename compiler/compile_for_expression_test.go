package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileForExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:  "defines the variable in the enclosing scope",
			input: "for x in [1]; end; x",
			expectedConstants: []any{
				1,
				":each",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpSetGlobal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpArray, 1),
				bytecode.Make(bytecode.OpCloseBlock, 2, 0),
				bytecode.Make(bytecode.OpSend, 1, 0, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name: "captures the variable and body locals of a method",
			input: `
				def m
					for x in [1]
						y = x
					end
					y
				end
			`,
			expectedConstants: []any{
				1,
				":each",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpSetFree, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetFree, 0),
					bytecode.Make(bytecode.OpSetFree, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
				":m",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpPushConstant, 0),
					bytecode.Make(bytecode.OpArray, 1),
					bytecode.Make(bytecode.OpCloseBlock, 2, 2),
					bytecode.Make(bytecode.OpSend, 1, 0, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
