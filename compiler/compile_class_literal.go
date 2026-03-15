package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	ast "emerald/parser/ast"
)

func (c *Compiler) compileClassLiteral(node *ast.ClassLiteral) {
	name := node.Name.Value

	// Emit a parent class
	// OpOpenClass expects this top be top of stack
	if node.Parent == nil {
		// If no parent is specified, it inherits from core.Object
		c.emitConstantGet(core.Object.Name, node.Token)
	} else {
		c.Compile(node.Parent)
	}

	c.emit(bytecode.OpOpenClass, node.Token, c.addConstant(core.NewSymbol(name)))

	c.compileStatementsWithReturnValue(node.Body.Statements, node.Body.Token)

	if c.lastInstructionIs(bytecode.OpPop) {
		c.replaceLastInstructionWith(bytecode.OpUnwrapContext)
	} else {
		c.emit(bytecode.OpUnwrapContext, node.Token)
	}
}

func (c *Compiler) compileStaticClassLiteral(node *ast.StaticClassLiteral) {
	c.emit(bytecode.OpStaticTrue, node.Token)

	c.Compile(node.Body)

	if c.lastInstructionIs(bytecode.OpPop) {
		lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
		c.replaceInstruction(lastPos, bytecode.Make(bytecode.OpStaticFalse))
		c.scopes[c.scopeIndex].lastInstruction.Opcode = bytecode.OpStaticFalse
	} else {
		c.emit(bytecode.OpStaticFalse, node.Token)
	}
}
