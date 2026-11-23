package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

// compileWhileExpression compiles `while` and `until` loops, which differ only in the jump that leaves the loop
func (c *Compiler) compileWhileExpression(node *ast.WhileExpression) {
	conditionPosition := len(c.currentInstructions())

	c.Compile(node.Condition)

	exitJump := bytecode.OpJumpNotTruthy
	if node.IsUntil() {
		exitJump = bytecode.OpJumpTruthy
	}

	// Emit the exit jump with a bogus value
	exitJumpPos := c.emit(exitJump, node.Token, 9999)
	c.emit(bytecode.OpPop, node.Token)

	scope := &c.scopes[c.scopeIndex]
	scope.loops = append(scope.loops, loop{conditionPosition: conditionPosition})
	c.Compile(node.Consequence)
	breakJumps := scope.loops[len(scope.loops)-1].breakJumps
	scope.loops = scope.loops[:len(scope.loops)-1]

	c.emit(bytecode.OpJump, node.Token, conditionPosition)

	afterConsequencePos := len(c.currentInstructions())
	c.changeOperand(exitJumpPos, afterConsequencePos)

	c.emit(bytecode.OpNull, node.Token)

	// `break` leaves its value on the stack in place of the OpNull above
	for _, pos := range breakJumps {
		c.changeOperand(pos, len(c.currentInstructions()))
	}
}
