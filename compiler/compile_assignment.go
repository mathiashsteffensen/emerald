package compiler

import (
	"emerald/ast"
	"emerald/core"
)

func (c *Compiler) compileAssignment(node *ast.AssignmentExpression) error {
	err := c.Compile(node.Value)
	if err != nil {
		return err
	}

	switch name := node.Name.(type) {
	case *ast.IdentifierExpression:
		symbol := c.symbolTable.Define(name.String())
		switch symbol.Scope {
		case GlobalScope:
			c.emit(OpSetGlobal, symbol.Index)
		case LocalScope:
			c.emit(OpSetLocal, symbol.Index)
		}
	case *ast.InstanceVariable:
		c.emit(OpInstanceVarSet, c.addConstant(core.NewSymbol(name.Value)))
	}

	return nil
}
