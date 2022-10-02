package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileStringLiteral(node *ast.StringLiteral) {
	str := core.NewString(node.Value)
	c.emit(OpPushConstant, c.addConstant(str))
}

func (c *Compiler) compileStringTemplate(node *ast.StringTemplate) error {
	beforeOpCount := c.opCount

	err := c.compileStringTemplateString(node.Chain)

	c.emit(OpStringJoin, c.opCount-beforeOpCount)

	return err
}

func (c *Compiler) compileStringTemplateString(node *ast.StringTemplateChainString) error {
	c.compileStringLiteral(node.StringLiteral)

	if node.Next != nil {
		return c.compileStringTemplateExpression(node.Next)
	}

	return nil
}

func (c *Compiler) compileStringTemplateExpression(node *ast.StringTemplateChainExpression) error {
	err := c.Compile(node.Expression)
	if err != nil {
		return err
	}

	if node.Next != nil {
		return c.compileStringTemplateString(node.Next)
	}

	return nil
}
