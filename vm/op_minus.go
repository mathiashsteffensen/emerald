package vm

import (
	"emerald/core"
)

func (vm *VM) executeOpMinus() {
	operand := vm.pop()

	vm.push(vm.Send(operand, "-@", core.NULL))
}
