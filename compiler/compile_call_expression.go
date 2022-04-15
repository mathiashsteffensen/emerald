package compiler

import (
	"emerald/ast"
	"emerald/core"
)

func (c *Compiler) compileCallExpression(node *ast.CallExpression) error {
	method := core.NewSymbol(node.Method.(*ast.IdentifierExpression).Value)

	c.emit(OpPushConstant, c.addConstant(method))

	for _, argument := range node.Arguments {
		err := c.Compile(argument)
		if err != nil {
			return err
		}
	}

	c.emit(OpSend, len(node.Arguments))

	return nil
}