package vm

import (
	"emerald/core"
	"emerald/object"
	"fmt"
)

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	typ := operand.ParentClass().(*object.Class).Name
	if typ != core.Integer.Name {
		return fmt.Errorf("unsupported type for negation: %s", typ)
	}

	value := operand.(*core.IntegerInstance).Value

	return vm.push(core.NewInteger(-value))
}
