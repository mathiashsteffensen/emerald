package vm

import (
	"emerald/core"
	"emerald/object"
)

func (vm *VM) executeOpMinus() {
	operand := vm.pop()

	vm.push(vm.Send(operand, "-@", core.NULL, map[string]object.EmeraldValue{}))
}
