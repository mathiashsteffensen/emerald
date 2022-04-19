package compiler

import "testing"

func TestCompileModuleLiteral(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "module with no methods",
			input:             "module MyMod; end",
			expectedConstants: []any{"module:MyMod"},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
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
			expectedConstants: []any{"module:MyMod", ":method", []Instructions{Make(OpReturn)}},
			expectedInstructions: []Instructions{
				Make(OpPushConstant, 0),
				Make(OpSetGlobal, 0),
				Make(OpOpenClass),
				Make(OpPushConstant, 1),
				Make(OpPushConstant, 2),
				Make(OpDefineMethod),
				Make(OpCloseClass),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
