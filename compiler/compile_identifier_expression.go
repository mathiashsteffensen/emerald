package compiler

import (
	"emerald/ast"
	"emerald/core"
)

func (c *Compiler) compileIdentifierExpression(node *ast.IdentifierExpression) {
	symbol, ok := c.symbolTable.Resolve(node.Value)
	if ok {
		c.emitSymbol(symbol)
	} else {
		c.emit(OpPushConstant, c.addConstant(core.NewSymbol(node.Value)))
		c.emit(OpNull)
		c.emit(OpSend)
	}
}
