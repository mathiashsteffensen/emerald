package compiler

import (
	"emerald/core"
	ast "emerald/parser/ast"
)

func (c *Compiler) compileClassLiteral(node *ast.ClassLiteral) {
	name := node.Name.Value

	// Emit a parent class
	// OpOpenClass expects this top be top of stack
	if node.Parent == nil {
		// If no parent is specified, it inherits from core.Object
		c.emitConstantGet(core.Object.Name)
	} else {
		c.Compile(node.Parent)
	}

	c.emit(OpOpenClass, c.addConstant(core.NewSymbol(name)))

	c.compileStatementsWithReturnValue(node.Body.Statements)

	if c.lastInstructionIs(OpPop) {
		c.replaceLastInstructionWith(OpUnwrapContext)
	} else {
		c.emit(OpUnwrapContext)
	}
}

func (c *Compiler) compileStaticClassLiteral(node *ast.StaticClassLiteral) {
	c.emit(OpStaticTrue)

	c.Compile(node.Body)

	if c.lastInstructionIs(OpPop) {
		lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
		c.replaceInstruction(lastPos, Make(OpStaticFalse))
		c.scopes[c.scopeIndex].lastInstruction.Opcode = OpStaticFalse
	} else {
		c.emit(OpStaticFalse)
	}
}
