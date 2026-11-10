package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type ForExpression struct {
	Token    lexer.Token // The 'for' token
	Variable IdentifierExpression
	Iterable Expression
	Body     *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, "for ")
	out.WriteString(fe.Variable.String(0))
	out.WriteString(" in ")
	out.WriteString(fe.Iterable.String(0))
	out.WriteString("\n")
	out.WriteString(fe.Body.String(indent + 1))
	indented(&out, indent, "end")

	return out.String()
}
