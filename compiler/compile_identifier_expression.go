package compiler

import (
	"emerald/ast"
	"emerald/core"
)

func (c *Compiler) compileIdentifierExpression(node ast.Expression) {
	switch node := node.(type) {
	case *ast.IdentifierExpression:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if ok {
			c.emitSymbol(symbol)
		} else {
			c.emit(OpPushConstant, c.addConstant(core.NewSymbol(node.Value)))
			c.emit(OpNull)
			c.emit(OpSend)
		}
	case *ast.InstanceVariable:
		c.emit(OpInstanceVarGet, c.addConstant(core.NewSymbol(node.Value)))
	}
}
