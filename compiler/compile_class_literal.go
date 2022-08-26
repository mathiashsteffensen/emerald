package compiler

import (
	"emerald/ast"
	"emerald/core"
	"emerald/object"
)

func (c *Compiler) compileClassLiteral(node *ast.ClassLiteral) error {
	name := node.Name.Value
	class := object.NewClass(name, core.Object, core.Object.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	c.emitConstantGetOrSet(name, class)
	c.emit(OpOpenClass)

	if len(node.Body.Statements) == 0 {
		c.emit(OpNull)
	} else {
		err := c.Compile(node.Body)
		if err != nil {
			return err
		}
	}

	if c.lastInstructionIs(OpPop) {
		lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
		c.replaceInstruction(lastPos, Make(OpCloseClass))
		c.scopes[c.scopeIndex].lastInstruction.Opcode = OpCloseClass
	} else {
		c.emit(OpCloseClass)
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
