package compiler

import "emerald/ast"

func (c *Compiler) compileMethodCall(node *ast.MethodCall) error {
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	c.emit(OpSetExecutionContext)

	err = c.compileCallExpression(node.CallExpression)
	if err != nil {
		return err
	}

	c.emit(OpResetExecutionContext)

	return nil
}
