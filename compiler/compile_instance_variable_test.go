package compiler

import "testing"

func TestCompileInstanceVariable(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:  "setting an instance var",
			input: "@var = 2 + 6",
			expectedConstants: []any{
				6,
				2,
				":@var",
			},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpPushConstant, 1),
				Make(OpAdd),
				Make(OpInstanceVarSet, 2),
				Make(OpPop),
			},
		},
		{
			name:  "getting an instance var",
			input: "other_var = @var",
			expectedConstants: []any{
				":@var",
			},
			expectedInstructions: []Instructions{
				Make(OpInstanceVarGet, 0),
				Make(OpSetGlobal, 0),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
