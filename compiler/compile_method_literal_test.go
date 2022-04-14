package compiler

import "testing"

func TestCompileMethodLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `def method
				return 5 + 10
			end`,
			expectedConstants: []interface{}{
				5,
				10,
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpPushConstant, 1),
					Make(OpAdd),
					Make(OpReturnValue),
				},
				":method",
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
			expectedConstants: []interface{}{
				5,
				10,
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpPushConstant, 1),
					Make(OpAdd),
					Make(OpReturnValue),
				},
				":method",
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
			expectedConstants: []interface{}{
				5,
				10,
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpPop),
					Make(OpPushConstant, 1),
					Make(OpReturnValue),
				},
				":method",
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
			expectedConstants: []interface{}{
				[]Instructions{
					Make(OpReturn),
				},
				":method",
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
