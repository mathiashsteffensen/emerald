package compiler

import (
	"emerald/ast"
	"emerald/core"
	"emerald/object"
)

func (c *Compiler) compileModuleLiteral(node *ast.ModuleLiteral) error {
	name := node.Name.Value
	class := object.NewModule(name, object.BuiltInMethodSet{}, object.BuiltInMethodSet{}, core.Module)

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
