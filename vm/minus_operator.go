package vm

import (
	"emerald/object"
	"fmt"
)

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	typ := operand.ParentClass().(*object.Class).Name
	if typ != object.Integer.Name {
		return fmt.Errorf("unsupported type for negation: %s", typ)
	}

	value := operand.(*object.IntegerInstance).Value

	return vm.push(object.NewInteger(-value))
}
