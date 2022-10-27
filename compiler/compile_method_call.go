package compiler

import (
	"emerald/parser/ast"
)

func (c *Compiler) compileMethodCall(node *ast.MethodCall) {
	c.Compile(node.Left)
	c.compileCallExpression(node.CallExpression)
}
