package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/object"
	"testing"
)

func TestCompileReturnStatement(t *testing.T) {
	tests := []struct {
		name                      string
		input                     string
		expectedBlockInstructions []bytecode.Instructions
	}{
		{
			name:                      "in a method",
			input:                     "def foo; return 1; end",
			expectedBlockInstructions: []bytecode.Instructions{bytecode.Make(bytecode.OpPushConstant, 0), bytecode.Make(bytecode.OpReturnValue)},
		},
		{
			name:                      "in a block",
			input:                     "[].each { return 1 }",
			expectedBlockInstructions: []bytecode.Instructions{bytecode.Make(bytecode.OpPushConstant, 1), bytecode.Make(bytecode.OpReturnFromBlock), bytecode.Make(bytecode.OpReturn)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, program := parse(tt.input)

			rt := core.NewRuntime()
			rt.Init()

			New(l, rt).Compile(program)

			for _, constant := range rt.Heap.ConstantPool {
				if block, ok := constant.Heap.(*object.Block); ok {
					if err := testInstructions(tt.expectedBlockInstructions, block.Instructions); err != nil {
						t.Error(err)
					}
				}
			}
		})
	}
}

func TestCompileReturnStatementMarksNonLocalReturn(t *testing.T) {
	l, program := parse("def outer; [].each { [].map { return 1 } }; end; def plain; [].each { next 1 }; end")

	rt := core.NewRuntime()
	rt.Init()

	New(l, rt).Compile(program)

	expected := map[string]bool{
		"outer":                     true,
		"block in outer":            true,
		"block (2 levels) in outer": true,
		"plain":                     false,
		"block in plain":            false,
	}

	for _, constant := range rt.Heap.ConstantPool {
		block, ok := constant.Heap.(*object.Block)
		if !ok {
			continue
		}

		want, ok := expected[block.Name]
		if !ok {
			t.Errorf("unexpected block %q", block.Name)
			continue
		}

		if block.HasNonLocalReturn != want {
			t.Errorf("%s: HasNonLocalReturn=%t, want=%t", block.Name, block.HasNonLocalReturn, want)
		}

		delete(expected, block.Name)
	}

	for name := range expected {
		t.Errorf("block %q was not compiled", name)
	}
}
