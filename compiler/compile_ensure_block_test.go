package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/object"
	"reflect"
	"testing"
)

func TestCompileEnsureBlock(t *testing.T) {
	tests := []struct {
		name                      string
		input                     string
		expectedBlockInstructions []bytecode.Instructions
		expectedExceptionTable    []object.ExceptionTableEntry
	}{
		{
			name:  "after the body",
			input: "def foo\n1\nensure\n2\nend",
			expectedBlockInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpReturnValue),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpEnsureEnd),
			},
			expectedExceptionTable: []object.ExceptionTableEntry{
				{StartIP: 0, EndIP: 3, HandlerIP: 4, Ensure: true},
			},
		},
		{
			name:  "after a rescue clause, which handles an error first",
			input: "def foo\n1\nrescue\n2\nensure\n3\nend",
			expectedBlockInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpPushConstant, 0),
				bytecode.Make(bytecode.OpReturnValue),
				bytecode.Make(bytecode.OpJump, 11),
				bytecode.Make(bytecode.OpPushConstant, 1),
				bytecode.Make(bytecode.OpReturnValue),
				bytecode.Make(bytecode.OpPushConstant, 2),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpEnsureEnd),
			},
			expectedExceptionTable: []object.ExceptionTableEntry{
				{StartIP: 0, EndIP: 4, HandlerIP: 7},
				{StartIP: 0, EndIP: 10, HandlerIP: 11, Ensure: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, program := parse(tt.input)

			rt := core.NewRuntime()
			rt.Init()

			New(l, rt).Compile(program)

			for _, constant := range rt.Heap.ConstantPool {
				block, ok := constant.Heap.(*object.Block)
				if !ok {
					continue
				}

				if err := testInstructions(tt.expectedBlockInstructions, block.Instructions); err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(block.ExceptionTable, tt.expectedExceptionTable) {
					t.Errorf("wrong exception table.\nwant=%+v\ngot=%+v", tt.expectedExceptionTable, block.ExceptionTable)
				}
			}
		})
	}
}
