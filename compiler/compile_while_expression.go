package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileWhileExpression(node *ast.WhileExpression) {
	conditionPosition := len(c.currentInstructions())

	c.Compile(node.Condition)

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, node.Token, 9999)
	c.emit(bytecode.OpPop, node.Token)

	c.Compile(node.Consequence)

	c.emit(bytecode.OpJump, node.Token, conditionPosition)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	c.emit(bytecode.OpNull, node.Token)
}
