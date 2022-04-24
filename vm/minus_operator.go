package vm

import (
	"emerald/core"
	"fmt"
)

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	switch operand := operand.(type) {
	case *core.IntegerInstance:
		return vm.push(core.NewInteger(-operand.Value))
	case *core.FloatInstance:
		return vm.push(core.NewFloat(-operand.Value))
	default:
		return fmt.Errorf("unsupported type for negation: %s", operand.Class().Super())
	}
}
