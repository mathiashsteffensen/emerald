package vm

import (
	"emerald/core"
	"emerald/object"
)

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	var result object.EmeraldValue

	if core.IsTruthy(operand) {
		result = core.FALSE
	} else {
		result = core.TRUE
	}

	return vm.push(result)
}
