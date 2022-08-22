package vm

import (
	"emerald/core"
	"fmt"
)

func (vm *VM) executeMinusOperator() {
	operand := vm.pop()

	switch operand := operand.(type) {
	case *core.IntegerInstance:
		vm.push(core.NewInteger(-operand.Value))
	case *core.FloatInstance:
		vm.push(core.NewFloat(-operand.Value))
	default:
		panic(fmt.Errorf("unsupported type for negation: %s", operand.Class().Super()))
	}
}
