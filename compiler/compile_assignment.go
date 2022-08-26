package compiler

import (
	"emerald/ast"
	"emerald/core"
	"unicode"
)

func (c *Compiler) compileAssignment(node *ast.AssignmentExpression) error {
	err := c.Compile(node.Value)
	if err != nil {
		return err
	}

	switch name := node.Name.(type) {
	case *ast.IdentifierExpression:
		if unicode.IsUpper(rune(name.Value[0])) {
			c.emit(OpConstantSet, c.addConstant(core.NewSymbol(name.Value)))
			return nil
		}

		stringName := name.String()

		symbol, ok := c.symbolTable.Resolve(stringName)
		if !ok {
			symbol = c.symbolTable.Define(stringName)
		}

		switch symbol.Scope {
		case GlobalScope:
			c.emit(OpSetGlobal, symbol.Index)
		case LocalScope:
			c.emit(OpSetLocal, symbol.Index)
		}
	case *ast.InstanceVariable:
		c.emit(OpInstanceVarSet, c.addConstant(core.NewSymbol(name.Value)))
	case *ast.GlobalVariable:
		symbol, ok := c.symbolTable.Resolve(name.Value)
		if !ok {
			symbol = c.symbolTable.DefineGlobal(name.Value)
		}

		c.emit(OpSetGlobal, symbol.Index)
	}

	return nil
}
