package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileIfExpression(node *ast.IfExpression) {
	c.Compile(node.Condition)

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, node.Token, 9999)
	c.emit(bytecode.OpPop, node.Token)

	if node.Consequence == nil {
		c.emit(bytecode.OpNull, node.Token)
	} else {
		c.compileStatementsWithReturnValue(node.Consequence.Statements, node.Token)

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}
	}

	// Emit an `OpJump` with a bogus value
	jumpPos := c.emit(bytecode.OpJump, node.Token, 9999)

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
		c.emit(bytecode.OpNull, node.Token)
	} else {
		c.compileStatementsWithReturnValue(node.Alternative.Statements, node.Token)
		if c.lastInstructionIs(bytecode.OpPop) {
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
	jumpNotTruthyPos := c.emit(bytecode.OpJumpNotTruthy, elsIf.Consequence.Token, 9999)
	c.emit(bytecode.OpPop, elsIf.Consequence.Token)

	if elsIf.Consequence == nil {
		c.emit(bytecode.OpNull, elsIf.Consequence.Token)
	} else {
		c.compileStatementsWithReturnValue(elsIf.Consequence.Statements, elsIf.Consequence.Token)

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}
	}

	// Emit an `OpJump` with a bogus value
	jumpPos := c.emit(bytecode.OpJump, elsIf.Consequence.Token, 9999)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	return jumpPos
}
