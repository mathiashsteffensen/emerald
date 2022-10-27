package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileModuleLiteral(node *ast.ModuleLiteral) {
	name := node.Name.Value

	c.emit(OpOpenModule, c.addConstant(core.NewSymbol(name)))

	c.compileStatementsWithReturnValue(node.Body.Statements)

	if c.lastInstructionIs(OpPop) {
		c.replaceLastInstructionWith(OpUnwrapContext)
	} else {
		c.emit(OpUnwrapContext)
	}
}
