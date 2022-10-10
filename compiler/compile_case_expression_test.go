package compiler

import "testing"

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
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpConstantGet, 1),
				Make(OpConstantGet, 2),
				Make(OpCheckCaseEqual, 2, 19),
				Make(OpPushConstant, 3),
				Make(OpJump, 23),
				Make(OpPop),
				Make(OpPushConstant, 4),
				Make(OpPop),
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
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpConstantGet, 1),
				Make(OpCheckCaseEqual, 1, 16),
				Make(OpPushConstant, 2),
				Make(OpJump, 33),
				Make(OpConstantGet, 3),
				Make(OpCheckCaseEqual, 1, 29),
				Make(OpPushConstant, 4),
				Make(OpJump, 33),
				Make(OpPop),
				Make(OpPushConstant, 5),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
