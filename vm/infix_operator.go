package vm

import (
	"emerald/compiler"
	"emerald/object"
)

// Operators are just methods so here we map opcodes to method names
var infixOperators = map[compiler.Opcode]string{
	compiler.OpAdd:         "+",
	compiler.OpSub:         "-",
	compiler.OpDiv:         "/",
	compiler.OpMul:         "*",
	compiler.OpLessThan:    "<",
	compiler.OpGreaterThan: ">",
	compiler.OpEqual:       "==",
	compiler.OpNotEqual:    "!=",
}

// execute the infix operator expression and push the result to the top of the stack
func (vm *VM) executeInfixOperator(op string, left, right object.EmeraldValue) error {
	result := left.SEND(func(block *object.Block, args ...object.EmeraldValue) object.EmeraldValue {
		return object.NULL
	}, op, left, nil, right)

	return vm.push(result)
}

// Decode bytecode into a left value, an operator and a right value
func (vm *VM) decodeInfixOperator(op compiler.Opcode) (object.EmeraldValue, string, object.EmeraldValue) {
	right := vm.pop()
	left := vm.pop()

	operator := infixOperators[op]

	return left, operator, right
}

func isInfixOperator(op compiler.Opcode) bool {
	_, ok := infixOperators[op]

	return ok
}
