package compiler

import (
	"emerald/parser/ast"
)

func (c *Compiler) compileIfExpression(node *ast.IfExpression) {
	c.Compile(node.Condition)

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)
	c.emit(OpPop)

	if node.Consequence == nil {
		c.emit(OpNull)
	} else {
		c.Compile(node.Consequence)

		if c.lastInstructionIs(OpPop) {
			c.removeLastPop()
		}
	}

	// Emit an `OpJump` with a bogus value
	jumpPos := c.emit(OpJump, 9999)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	jumpPositions := []int{jumpPos}

	if node.ElseIfs != nil {
		for _, elseIf := range node.ElseIfs {
			elseIfPosition := c.compileElsifBranch(elseIf)
			jumpPositions = append(jumpPositions, elseIfPosition)
		}
	}

	if node.Alternative == nil {
		c.emit(OpNull)
	} else {
		c.Compile(node.Alternative)
		if c.lastInstructionIs(OpPop) {
			c.removeLastPop()
		}
	}

	afterAlternativePos := len(c.currentInstructions())
	for _, position := range jumpPositions {
		c.changeOperand(position, afterAlternativePos)
	}
}

func (c *Compiler) compileElsifBranch(elsIf ast.ElseIf) int {
	c.Compile(elsIf.Condition)

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)
	c.emit(OpPop)

	if elsIf.Consequence == nil {
		c.emit(OpNull)
	} else {
		c.Compile(elsIf.Consequence)

		if c.lastInstructionIs(OpPop) {
			c.removeLastPop()
		}
	}

	// Emit an `OpJump` with a bogus value
	jumpPos := c.emit(OpJump, 9999)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	return jumpPos
}
