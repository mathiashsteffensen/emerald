package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type RescueBlock struct {
	Token              lexer.Token // The 'rescue' token
	Body               *BlockStatement
	CaughtErrorClasses []Expression
	ErrorVarName       Expression
}

func (rb *RescueBlock) TokenLiteral() string { return rb.Token.Literal }
func (rb *RescueBlock) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, rb.TokenLiteral())
	out.WriteString("\n")
	out.WriteString(rb.Body.String(indent + 1))

	return out.String()
}
