package vm

import "emerald/object"

// executeOpBlock pushes the block passed to the current call as a Proc, or nil
// when the call did not receive one. Native blocks, such as the wrappers
// Enumerable passes to each, are captured too.
func (vm *VM) executeOpBlock() {
	block := vm.ctx.Block
	if block.IsNil() {
		vm.push(vm.rt.NULL)
		return
	}

	switch block.Heap.(type) {
	case *object.ClosedBlock, *object.WrappedBuiltInMethod:
		vm.push(vm.rt.NewProc(block))
	default:
		vm.push(vm.rt.NULL)
	}
}
