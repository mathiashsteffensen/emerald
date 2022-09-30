package vm

import "emerald/core"

func (vm *VM) executeOpReturn() {
	vm.currentFiber().popFrame()

	vm.push(core.NULL)
}

func (vm *VM) executeOpReturnValue() {
	returnValue := vm.stack()[vm.currentFiber().sp-1]

	vm.currentFiber().popFrame()

	vm.push(returnValue)
}
