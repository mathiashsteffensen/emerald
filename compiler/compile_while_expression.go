package compiler

import "emerald/parser/ast"

func (c *Compiler) compileWhileExpression(node *ast.WhileExpression) error {
	conditionPosition := len(c.currentInstructions())

	err := c.Compile(node.Condition)
	if err != nil {
		return err
	}

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)
	c.emit(OpPop)

	err = c.Compile(node.Consequence)
	if err != nil {
		return err
	}

	c.emit(OpJump, conditionPosition)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	c.emit(OpNull)

	return nil
}
