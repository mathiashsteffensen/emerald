package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/heap"
	ast "emerald/parser/ast"
	"unicode"
)

func (c *Compiler) compileAssignment(node *ast.AssignmentExpression) {
	c.Compile(node.Value)

	switch name := node.Name.(type) {
	case ast.IdentifierExpression:
		if unicode.IsUpper(rune(name.Value[0])) {
			c.emit(bytecode.OpConstantSet, node.Token, c.addConstant(core.NewSymbol(name.Value)))
			return
		}

		stringName := name.String(0)

		symbol, ok := c.symbolTable.Resolve(stringName)
		if !ok {
			symbol = c.symbolTable.Define(stringName)
		}

		switch symbol.Scope {
		case heap.GlobalScope:
			c.emit(bytecode.OpSetGlobal, node.Token, symbol.Index)
		case heap.LocalScope:
			c.emit(bytecode.OpSetLocal, node.Token, symbol.Index)
		}
	case *ast.InstanceVariable:
		c.emit(bytecode.OpInstanceVarSet, node.Token, c.addConstant(core.NewSymbol(name.Value)))
	case *ast.GlobalVariable:
		symbol, ok := c.symbolTable.Resolve(name.Value)
		if !ok {
			symbol = c.symbolTable.DefineGlobal(name.Value)
		}

		c.emit(bytecode.OpSetGlobal, node.Token, symbol.Index)
	}
}
