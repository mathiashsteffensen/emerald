package compiler

import "emerald/ast"

func (c *Compiler) compileGlobalAssignment(node *ast.AssignmentExpression) error {
	err := c.Compile(node.Value)
	if err != nil {
		return err
	}

	symbol := c.symbolTable.Define(node.Name.String())

	c.emit(OpSetGlobal, symbol.Index)
	c.emit(OpGetGlobal, symbol.Index)

	return nil
}
