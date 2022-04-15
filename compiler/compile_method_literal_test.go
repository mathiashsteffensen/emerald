package compiler

import "testing"

func TestCompileMethodLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `def method
				return 5 + 10
			end`,
			expectedConstants: []interface{}{
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
			expectedConstants: []interface{}{
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
			expectedConstants: []interface{}{
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
			expectedConstants: []interface{}{
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
	}
	runCompilerTests(t, tests)
}
