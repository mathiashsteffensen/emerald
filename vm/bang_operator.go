package vm

import "emerald/object"

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	var result object.EmeraldValue

	if isTruthy(operand) {
		result = object.FALSE
	} else {
		result = object.TRUE
	}

	return vm.push(result)
}
