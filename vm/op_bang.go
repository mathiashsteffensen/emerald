package vm

import (
	"emerald/object"
)

func (vm *VM) executeBangOperator() {
	operand := vm.pop()

	result := vm.Send(operand, "!@", vm.rt.NULL, map[string]object.EmeraldValue{})

	if !vm.ExceptionIsRaised() {
		vm.push(result)
	}
}
