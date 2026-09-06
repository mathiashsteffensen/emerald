package vm

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/object"
	"testing"
)

func TestLegacyClosedBlockFreeVariables(t *testing.T) {
	for _, direct := range []bool{false, true} {
		name := "constructor"
		if direct {
			name = "struct literal"
		}
		t.Run(name, func(t *testing.T) {
			rt := core.NewRuntime()
			rt.Init()
			free := []object.EmeraldValue{rt.NewInteger(7)}
			child := &object.Block{FreeBindings: []object.FreeBinding{{Index: 0}}}
			childIndex := rt.Heap.AddConstant(object.NewHeapObject(child))
			valueIndex := rt.Heap.AddConstant(rt.NewInteger(9))
			var instructions bytecode.Instructions
			for _, instruction := range []bytecode.Instructions{
				bytecode.Make(bytecode.OpGetFree, 0),
				bytecode.Make(bytecode.OpPushConstant, valueIndex),
				bytecode.Make(bytecode.OpSetFree, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpGetFree, 0),
				bytecode.Make(bytecode.OpCloseBlock, childIndex, 1),
			} {
				instructions = append(instructions, instruction...)
			}
			block := &object.Block{Bytecode: bytecode.Bytecode{Instructions: instructions}}
			closed := object.NewClosedBlock(nil, block, free, "legacy.em", object.PUBLIC)
			if direct {
				closed = &object.ClosedBlock{Block: block, FreeVariables: free}
			}
			vm := New("legacy.em", &block.Bytecode, rt)
			vm.currentFiber().currentFrame().block = closed
			closed.FreeVariables[0] = rt.NewInteger(8)
			vm.Run()
			nested := vm.pop().Heap.(*object.ClosedBlock)
			if vm.pop() != rt.NewInteger(9) || vm.pop() != rt.NewInteger(8) || free[0] != rt.NewInteger(9) {
				t.Fatal("GetFree/SetFree did not use the legacy values")
			}
			if nested.FreeVariable(0) != &free[0] {
				t.Fatal("nested closure did not capture the legacy variable's cell")
			}
			*nested.FreeVariable(0) = rt.NewInteger(10)
			if closed.FreeVariables[0] != rt.NewInteger(10) {
				t.Fatal("nested mutation did not update the legacy field")
			}
		})
	}
}
