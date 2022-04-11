package compiler

import "emerald/ast"

func (c *Compiler) compileIfExpression(node *ast.IfExpression) error {
	err := c.Compile(node.Condition)
	if err != nil {
		return err
	}

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)

	err = c.Compile(node.Consequence)
	if err != nil {
		return err
	}

	if c.lastInstructionIsPop() {
		c.removeLastPop()
	}

	// Emit an `OpJump` with a bogus value
	jumpPos := c.emit(OpJump, 9999)

	afterConsequencePos := len(c.instructions)
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	if node.Alternative == nil {
		c.emit(OpNull)
	} else {
		err := c.Compile(node.Alternative)
		if err != nil {
			return err
		}
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}
	}

	afterAlternativePos := len(c.instructions)
	c.changeOperand(jumpPos, afterAlternativePos)

	return nil
}
