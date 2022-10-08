package compiler

import (
	"emerald/parser/ast"
)

func (c *Compiler) compileIfExpression(node *ast.IfExpression) error {
	err := c.Compile(node.Condition)
	if err != nil {
		return err
	}

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)
	c.emit(OpPop)

	if node.Consequence == nil {
		c.emit(OpNull)
	} else {
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

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
			elseIfPosition, err := c.compileElsifBranch(elseIf)
			if err != nil {
				return err
			}
			jumpPositions = append(jumpPositions, elseIfPosition)
		}
	}

	if node.Alternative == nil {
		c.emit(OpNull)
	} else {
		err := c.Compile(node.Alternative)
		if err != nil {
			return err
		}
		if c.lastInstructionIs(OpPop) {
			c.removeLastPop()
		}
	}

	afterAlternativePos := len(c.currentInstructions())
	for _, position := range jumpPositions {
		c.changeOperand(position, afterAlternativePos)
	}

	return nil
}

func (c *Compiler) compileElsifBranch(elsIf ast.ElseIf) (int, error) {
	err := c.Compile(elsIf.Condition)
	if err != nil {
		return 0, err
	}

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)
	c.emit(OpPop)

	if elsIf.Consequence == nil {
		c.emit(OpNull)
	} else {
		err = c.Compile(elsIf.Consequence)
		if err != nil {
			return 0, err
		}

		if c.lastInstructionIs(OpPop) {
			c.removeLastPop()
		}
	}

	// Emit an `OpJump` with a bogus value
	jumpPos := c.emit(OpJump, 9999)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	return jumpPos, nil
}
