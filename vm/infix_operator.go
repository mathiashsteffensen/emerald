package vm

import (
	"emerald/compiler"
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
