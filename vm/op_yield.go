package vm

import "emerald/compiler"

func (vm *VM) executeOpYield(ins compiler.Instructions, ip int) {
	numArgs := vm.readUint8(ins, ip)
	args := vm.stack[vm.sp-int(numArgs) : vm.sp]

	vm.push(vm.ctx.Yield(args...))
}
