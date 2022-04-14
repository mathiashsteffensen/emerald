package compiler

import "emerald/ast"

func (c *Compiler) compileAssignment(node *ast.AssignmentExpression) error {
	err := c.Compile(node.Value)
	if err != nil {
		return err
	}

	symbol := c.symbolTable.Define(node.Name.String())

	switch symbol.Scope {
	case GlobalScope:
		c.emit(OpSetGlobal, symbol.Index)
		c.emit(OpGetGlobal, symbol.Index)
	case LocalScope:
		c.emit(OpSetLocal, symbol.Index)
		c.emit(OpGetLocal, symbol.Index)
	}

	return nil
}
