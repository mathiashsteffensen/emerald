package compiler

import (
	"emerald/bytecode"
	"testing"
)

func TestCompileModuleLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "module with no methods",
			input:             "module MyMod; end",
			expectedConstants: []any{":MyMod"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpOpenModule, 0),
				bytecode.Make(bytecode.OpNull),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			name: "module with method",
			input: `module MyMod
				def method; end
			end`,
			expectedConstants: []any{":MyMod", ":method", []bytecode.Instructions{bytecode.Make(bytecode.OpReturn)}},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpOpenModule, 0),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpDefineMethod),
				bytecode.Make(bytecode.OpUnwrapContext),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
