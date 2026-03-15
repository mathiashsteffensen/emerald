package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

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

		lastOpCheckCaseEqualPosition = c.emit(bytecode.OpCheckCaseEqual, clause.Token, lastOpCheckCaseEqualMatchersLength, 9999)
		lastOpCheckCaseEqualMatchersLength = len(clause.Matchers)

		c.Compile(clause.Consequence)

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}

		// Emit an OpJump with a bogus position, position will be set to right after else clause
		// when the else clause has been compiled
		opJumpPositions = append(opJumpPositions, c.emit(bytecode.OpJump, clause.Token, 9998))
	}

	c.changeOperand(lastOpCheckCaseEqualPosition, lastOpCheckCaseEqualMatchersLength, len(c.currentInstructions()))

	c.emit(bytecode.OpPop, node.Token)

	c.compileStatementsWithReturnValue(node.Alternative.Statements, node.Alternative.Token)

	if c.lastInstructionIs(bytecode.OpPop) {
		c.removeLastPop()
	}

	for _, position := range opJumpPositions {
		c.changeOperand(position, len(c.currentInstructions()))
	}

	c.emit(bytecode.OpPop, node.Token)
}
