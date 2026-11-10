package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileReturnStatement(node *ast.ReturnStatement) {
	c.Compile(node.ReturnValue)
	if c.aborted {
		return
	}

	scope := &c.scopes[c.scopeIndex]

	// Inside a block, `return` leaves the method the block literal is written in
	if scope.block {
		scope.nonLocalReturn = true
		c.emit(bytecode.OpReturnFromBlock, node.Token)
		return
	}

	c.emit(bytecode.OpReturnValue, node.Token)
}
