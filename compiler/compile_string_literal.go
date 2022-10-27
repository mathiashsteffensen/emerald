package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileStringLiteral(node *ast.StringLiteral) {
	str := core.NewString(node.Value)
	c.emit(OpPushConstant, c.addConstant(str))
}

func (c *Compiler) compileStringTemplate(node *ast.StringTemplate) {
	c.compileStringTemplateString(node.Chain)

	c.emit(OpStringJoin, node.Count())
}

func (c *Compiler) compileStringTemplateString(node *ast.StringTemplateChainString) {
	c.compileStringLiteral(node.StringLiteral)

	if node.Next != nil {
		c.compileStringTemplateExpression(node.Next)
	}
}

func (c *Compiler) compileStringTemplateExpression(node *ast.StringTemplateChainExpression) {
	c.Compile(node.Expression)

	if node.Next != nil {
		c.compileStringTemplateString(node.Next)
	}
}
