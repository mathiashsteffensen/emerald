package compiler

import (
	"emerald/bytecode"
	ast "emerald/parser/ast"
	"unicode"
)

func (c *Compiler) compileIdentifierExpression(node ast.Expression) {
	switch node := node.(type) {
	case ast.IdentifierExpression:
		if unicode.IsUpper(rune(node.Value[0])) {
			// Constant reference
			c.emitConstantGet(node.Value, node.Token)
		} else {
			symbol, ok := c.symbolTable.Resolve(node.Value)
			if ok {
				// Variable reference
				c.emitSymbol(symbol, node.Token)
			} else {
				// Method call with no arguments
				methodIndex := c.addConstant(c.rt.NewSymbol(node.Value))
				c.emit(bytecode.OpSelf, node.Token) // Call on self
				c.emit(bytecode.OpNull, node.Token) // No block
				c.emit(bytecode.OpSend, node.Token, methodIndex, 0, 0)
			}
		}
	case *ast.InstanceVariable:
		c.emit(bytecode.OpInstanceVarGet, node.Token, c.addConstant(c.rt.NewSymbol(node.Value)))
	case *ast.GlobalVariable:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			symbol = c.symbolTable.DefineGlobal(node.Value)
		}
		c.emitSymbol(symbol, node.Token)
	}
}
