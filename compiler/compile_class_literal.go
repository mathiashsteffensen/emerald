package compiler

import (
	"emerald/core"
	ast "emerald/parser/ast"
)

func (c *Compiler) compileClassLiteral(node *ast.ClassLiteral) error {
	name := node.Name.Value

	// Emit a parent class
	// OpOpenClass expects this top be top of stack
	if node.Parent == nil {
		// If no parent is specified, it inherits from core.Object
		c.emitConstantGet(core.Object.Name)
	} else {
		err := c.Compile(node.Parent)
		if err != nil {
			return err
		}
	}

	c.emit(OpOpenClass, c.addConstant(core.NewSymbol(name)))

	err := c.compileStatementsWithReturnValue(node.Body.Statements)
	if err != nil {
		return err
	}

	if c.lastInstructionIs(OpPop) {
		c.replaceLastInstructionWith(OpUnwrapContext)
	} else {
		c.emit(OpUnwrapContext)
	}

	return nil
}

func (c *Compiler) compileStaticClassLiteral(node *ast.StaticClassLiteral) error {
	c.emit(OpStaticTrue)

	err := c.Compile(node.Body)
	if err != nil {
		return err
	}

	if c.lastInstructionIs(OpPop) {
		lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
		c.replaceInstruction(lastPos, Make(OpStaticFalse))
		c.scopes[c.scopeIndex].lastInstruction.Opcode = OpStaticFalse
	} else {
		c.emit(OpStaticFalse)
	}

	return nil
}
