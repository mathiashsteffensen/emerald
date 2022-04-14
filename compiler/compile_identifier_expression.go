package compiler

import (
	"emerald/ast"
	"emerald/object"
)

func (c *Compiler) compileIdentifierExpression(node *ast.IdentifierExpression) {
	symbol, ok := c.symbolTable.Resolve(node.Value)
	if ok {
		switch symbol.Scope {
		case GlobalScope:
			c.emit(OpGetGlobal, symbol.Index)
		case LocalScope:
			c.emit(OpGetLocal, symbol.Index)
		}
	} else {
		c.emit(OpPushConstant, c.addConstant(object.NewSymbol(node.Value)))
		c.emit(OpSend)
	}
}
