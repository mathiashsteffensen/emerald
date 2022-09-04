package vm

import "emerald/compiler"

func (vm *VM) executeOpYield(ins compiler.Instructions, ip int) {
	numArgs := vm.readUint8(ins, ip)
	args := vm.stack()[vm.currentFiber().sp-int(numArgs) : vm.currentFiber().sp]

	vm.push(vm.ctx.Yield(args...))
}
