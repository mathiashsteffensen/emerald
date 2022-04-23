package compiler

import (
	"emerald/ast"
	"emerald/object"
)

func (c *Compiler) compileModuleLiteral(node *ast.ModuleLiteral) error {
	name := node.Name.Value

	var (
		symbol Symbol
		ok     bool
	)
	symbol, ok = c.symbolTable.Resolve(name)
	if ok {
		c.emit(OpGetGlobal, symbol.Index)
	} else {
		symbol = c.symbolTable.Define(name)
		class := object.NewModule(name, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

		c.emit(OpPushConstant, c.addConstant(class.StaticModule))
		c.emit(OpSetGlobal, symbol.Index)
	}

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
