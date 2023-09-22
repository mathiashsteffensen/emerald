package vm

import (
	"emerald/compiler"
	"emerald/object"
)

func (vm *VM) executeOpYield(ins compiler.Instructions, ip int) {
	numArgs := vm.readUint8(ins, ip)
	args := vm.stack()[vm.currentFiber().sp-int(numArgs) : vm.currentFiber().sp]

	result := vm.ctx.Yield(map[string]object.EmeraldValue{}, args...)
	if vm.currentFiber().inRescue || !vm.ExceptionIsRaised() {
		vm.push(result)
	}
}
