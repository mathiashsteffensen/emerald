package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileCaseExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				case 2
				when Integer, String
					4
				else
					5
				end
			`,
			expectedConstants: []any{
				2,
				":Integer",
				":String",
				4,
				5,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpConstantGet, 1),
				bytecode.Make(bytecode.OpConstantGet, 2),
				bytecode.Make(bytecode.OpCheckCaseEqual, 2, 19),
				bytecode.Make(bytecode.OpPushConstant, 3),
				bytecode.Make(bytecode.OpJump, 23),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
				case 2
				when Integer
					3
				when String
					4
				else
					5
				end
			`,
			expectedConstants: []any{
				2,
				":Integer",
				3,
				":String",
				4,
				5,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpConstantGet, 1),
				bytecode.Make(bytecode.OpCheckCaseEqual, 1, 16),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpJump, 33),
				bytecode.Make(bytecode.OpConstantGet, 3),
				bytecode.Make(bytecode.OpCheckCaseEqual, 1, 29),
				bytecode.Make(bytecode.OpPushConstant, 4),
				bytecode.Make(bytecode.OpJump, 33),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpPushConstant, 5),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
