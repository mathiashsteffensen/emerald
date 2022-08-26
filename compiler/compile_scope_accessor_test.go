package compiler

import "testing"

func TestCompileScopeAccessor(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "MyMod::MyClass",
			expectedConstants: []any{
				":MyMod",
				":MyClass",
			},
			expectedInstructions: []Instructions{
				Make(OpConstantGet, 0),
				Make(OpScopedConstantGet, 1),
				Make(OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
