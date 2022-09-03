package compiler

import "testing"

func TestCompileModuleLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "module with no methods",
			input:             "module MyMod; end",
			expectedConstants: []any{":MyMod"},
			expectedInstructions: []Instructions{
				Make(OpOpenModule, 0),
				Make(OpNull),
				Make(OpUnwrapContext),
				Make(OpPop),
			},
		},
		{
			name: "module with method",
			input: `module MyMod
				def method; end
			end`,
			expectedConstants: []any{":MyMod", ":method", []Instructions{Make(OpReturn)}},
			expectedInstructions: []Instructions{
				Make(OpOpenModule, 0),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpDefineMethod),
				Make(OpUnwrapContext),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
