package compiler

import "emerald/parser/ast"

func (c *Compiler) compileCaseExpression(node *ast.CaseExpression) {
	c.Compile(node.Subject)

	var (
		lastOpCheckCaseEqualPosition       = -1
		lastOpCheckCaseEqualMatchersLength int
		opJumpPositions                    []int
	)

	for _, clause := range node.WhenClauses {
		if lastOpCheckCaseEqualPosition != -1 {
			c.changeOperand(lastOpCheckCaseEqualPosition, lastOpCheckCaseEqualMatchersLength, len(c.currentInstructions()))
		}

		for _, matcher := range clause.Matchers {
			c.Compile(matcher)
		}

		lastOpCheckCaseEqualPosition = c.emit(OpCheckCaseEqual, lastOpCheckCaseEqualMatchersLength, 9999)
		lastOpCheckCaseEqualMatchersLength = len(clause.Matchers)

		c.Compile(clause.Consequence)

		if c.lastInstructionIs(OpPop) {
			c.removeLastPop()
		}

		// Emit an OpJump with a bogus position, position will be set to right after else clause
		// when the else clause has been compiled
		opJumpPositions = append(opJumpPositions, c.emit(OpJump, 9998))
	}

	c.changeOperand(lastOpCheckCaseEqualPosition, lastOpCheckCaseEqualMatchersLength, len(c.currentInstructions()))

	c.emit(OpPop)

	c.compileStatementsWithReturnValue(node.Alternative.Statements)

	if c.lastInstructionIs(OpPop) {
		c.removeLastPop()
	}

	for _, position := range opJumpPositions {
		c.changeOperand(position, len(c.currentInstructions()))
	}

	c.emit(OpPop)
}
