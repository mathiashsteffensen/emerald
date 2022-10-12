package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type ClassLiteral struct {
	Token  lexer.Token // The class token
	Name   IdentifierExpression
	Parent Expression
	Body   *BlockStatement
}

func (cl *ClassLiteral) expressionNode()      {}
func (cl *ClassLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *ClassLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, "class ")
	out.WriteString(cl.Name.String(0))
	out.WriteString("\n")
	out.WriteString(cl.Body.String(indent + 1))
	indented(&out, indent, "end")

	return out.String()
}
