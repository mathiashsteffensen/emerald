package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileModuleLiteral(node *ast.ModuleLiteral) {
	name := node.Name.Value

	c.emit(bytecode.OpOpenModule, node.Token, c.addConstant(c.rt.NewSymbol(name)))

	c.compileStatementsWithReturnValue(node.Body.Statements, node.Body.Token)

	if c.lastInstructionIs(bytecode.OpPop) {
		c.replaceLastInstructionWith(bytecode.OpUnwrapContext)
	} else {
		c.emit(bytecode.OpUnwrapContext, node.Token)
	}
}
