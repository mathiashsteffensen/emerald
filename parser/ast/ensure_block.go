package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type EnsureBlock struct {
	Token lexer.Token // The 'ensure' token
	Body  *BlockStatement
}

func (rb *EnsureBlock) TokenLiteral() string { return rb.Token.Literal }
func (rb *EnsureBlock) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, "ensure\n")
	out.WriteString(rb.Body.String(indent + 1))

	return out.String()
}
