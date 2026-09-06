package vm

import (
	"emerald/bytecode"
	"emerald/object"
)

func (vm *VM) executeOpYield(ins bytecode.Instructions, ip int) {
	numArgs := vm.readUint8(ins, ip)
	argsStart := vm.currentFiber().sp - int(numArgs)
	args := vm.stack()[argsStart:vm.currentFiber().sp]

	result := vm.ctx.Yield(map[string]object.EmeraldValue{}, args...)
	if vm.executionError == nil && !vm.ExceptionIsRaised() {
		vm.currentFiber().sp = argsStart
		vm.push(result)
	}
}
