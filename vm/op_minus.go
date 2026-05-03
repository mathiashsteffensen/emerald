package vm

import (
	"emerald/object"
)

func (vm *VM) executeOpMinus() {
	operand := vm.pop()

	vm.push(vm.Send(operand, "-@", vm.rt.NULL, map[string]object.EmeraldValue{}))
}
