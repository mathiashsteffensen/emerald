package compiler

import (
	"emerald/core"
	ast "emerald/parser/ast"
	"unicode"
)

func (c *Compiler) compileIdentifierExpression(node ast.Expression) {
	switch node := node.(type) {
	case ast.IdentifierExpression:
		if unicode.IsUpper(rune(node.Value[0])) {
			c.emitConstantGet(node.Value)
		} else {
			symbol, ok := c.symbolTable.Resolve(node.Value)
			if ok {
				c.emitSymbol(symbol)
			} else {
				c.emit(OpSelf)
				c.emit(OpPushConstant, c.addConstant(core.NewSymbol(node.Value)))
				c.emit(OpNull)
				c.emit(OpSend)
			}
		}
	case *ast.InstanceVariable:
		c.emit(OpInstanceVarGet, c.addConstant(core.NewSymbol(node.Value)))
	case *ast.GlobalVariable:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			symbol = c.symbolTable.DefineGlobal(node.Value)
		}
		c.emitSymbol(symbol)
	}
}
