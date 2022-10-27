package compiler

import "emerald/parser/ast"

func (c *Compiler) compileWhileExpression(node *ast.WhileExpression) {
	conditionPosition := len(c.currentInstructions())

	c.Compile(node.Condition)

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)
	c.emit(OpPop)

	c.Compile(node.Consequence)

	c.emit(OpJump, conditionPosition)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	c.emit(OpNull)
}
