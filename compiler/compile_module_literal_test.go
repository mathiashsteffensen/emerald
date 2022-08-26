package compiler

import "testing"

func TestCompileModuleLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "module with no methods",
			input:             "module MyMod; end",
			expectedConstants: []any{":MyMod", "module:MyMod"},
			expectedInstructions: []Instructions{
				Make(OpConstantGetOrSet, 0, 1),
				Make(OpOpenClass),
				Make(OpNull),
				Make(OpCloseClass),
				Make(OpPop),
			},
		},
		{
			name: "module with method",
			input: `module MyMod
				def method; end
			end`,
			expectedConstants: []any{":MyMod", "module:MyMod", ":method", []Instructions{Make(OpReturn)}},
			expectedInstructions: []Instructions{
				Make(OpConstantGetOrSet, 0, 1),
				Make(OpOpenClass),
				Make(OpPushConstant, 2),
				Make(OpPushConstant, 3),
				Make(OpDefineMethod),
				Make(OpCloseClass),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
