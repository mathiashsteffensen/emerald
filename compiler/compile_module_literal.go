package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileModuleLiteral(node *ast.ModuleLiteral) error {
	name := node.Name.Value

	c.emit(OpOpenModule, c.addConstant(core.NewSymbol(name)))

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
