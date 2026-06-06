package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileCallExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			def one_arg(a); a; end
			one_arg(24)
			`,
			expectedConstants: []any{
				":one_arg",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
				":one_arg",
				24,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpSend, 2, 1, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: "call { |n| n + 2 } ",
			expectedConstants: []any{
				":call",
				2,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpCloseBlock, 2),
				bytecode.Make(bytecode.OpSend, 0, 0, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			def many_arg(a, b, c); a; b; c; end
			many_arg(24, 25, 26)
			`,
			expectedConstants: []any{
				":many_arg",
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 2),
					bytecode.Make(bytecode.OpReturnValue),
				},
				":many_arg",
				24,
				25,
				26,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpSelf),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpPushConstant, 5),
				bytecode.Make(bytecode.OpSend, 2, 3, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
