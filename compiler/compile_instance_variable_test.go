package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileInstanceVariable(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:  "setting an instance var",
			input: "@var = 2 + 6",
			expectedConstants: []any{
				2,
				6,
				":@var",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpInstanceVarSet, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name:  "getting an instance var",
			input: "other_var = @var",
			expectedConstants: []any{
				":@var",
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpInstanceVarGet, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
