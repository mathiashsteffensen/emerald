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

	scope := &c.scopes[c.scopeIndex]
	scope.loops = append(scope.loops, conditionPosition)
	c.Compile(node.Consequence)
	scope.loops = scope.loops[:len(scope.loops)-1]

	c.emit(bytecode.OpJump, node.Token, conditionPosition)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	c.emit(bytecode.OpNull, node.Token)
}
