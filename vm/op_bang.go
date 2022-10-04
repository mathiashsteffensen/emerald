package vm

import (
	"emerald/core"
)

func (vm *VM) executeBangOperator() {
	operand := vm.pop()

	vm.push(vm.Send(operand, "!@", core.NULL))
}
