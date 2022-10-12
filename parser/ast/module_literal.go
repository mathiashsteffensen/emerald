package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type ModuleLiteral struct {
	Token lexer.Token // The module token
	Body  *BlockStatement
	Name  IdentifierExpression
}

func (cl *ModuleLiteral) expressionNode()      {}
func (cl *ModuleLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *ModuleLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, "module ")
	out.WriteString(cl.Name.String(0))
	out.WriteString("\n")

	for _, value := range cl.Body.Statements {
		out.WriteString(value.String(indent+1) + "\n")
	}

	indented(&out, indent, "end")

	return out.String()
}
