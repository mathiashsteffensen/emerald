package vm

import "emerald/compiler"

func (vm *VM) executeOpYield(ins compiler.Instructions, ip int) {
	numArgs := vm.readUint8(ins, ip)
	args := vm.stack()[vm.currentFiber().sp-int(numArgs) : vm.currentFiber().sp]

	result := vm.ctx.Yield(args...)
	if vm.currentFiber().inRescue || !vm.ExceptionIsRaised() {
		vm.push(result)
	}
}
