package compiler

import "testing"

func TestAssignments(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			one = 1
			two = 2
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpGetGlobal, 0),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpSetGlobal, 1),
				Make(OpGetGlobal, 1),
				Make(OpPop),
			},
		},
		{
			input: `
			one = 1
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpGetGlobal, 0),
				Make(OpPop),
			},
		},
		{
			input: `
			one = 1
			two = one
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpGetGlobal, 0),
				Make(OpPop),
				Make(OpGetGlobal, 0),
				Make(OpSetGlobal, 1),
				Make(OpGetGlobal, 1),
				Make(OpPop),
			},
		},
		{
			input: `
			two = one = 1
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpGetGlobal, 0),
				Make(OpSetGlobal, 1),
				Make(OpGetGlobal, 1),
				Make(OpPop),
			},
		},
		{
			input: `
			num = 55
			def method; num; end
			`,
			expectedConstants: []interface{}{
				55,
				[]Instructions{
					Make(OpGetGlobal, 0),
					Make(OpReturnValue),
				},
				":method",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpGetGlobal, 0),
				Make(OpPop),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			input: `
			def method
				num = 55
			end
			`,
			expectedConstants: []interface{}{
				55,
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpSetLocal, 0),
					Make(OpGetLocal, 0),
					Make(OpReturnValue),
				},
				":method",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpDefineMethod),
				Make(OpPop),
			},
		},
		{
			input: `
			def method
				a = 55
				b = 77
				a + b
			end
			`,
			expectedConstants: []interface{}{
				55,
				77,
				[]Instructions{
					Make(OpPushConstant, 0),
					Make(OpSetLocal, 0),
					Make(OpGetLocal, 0),
					Make(OpPop),
					Make(OpPushConstant, 1),
					Make(OpSetLocal, 1),
					Make(OpGetLocal, 1),
					Make(OpPop),
					Make(OpGetLocal, 0),
					Make(OpGetLocal, 1),
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
	}
	runCompilerTests(t, tests)
}
