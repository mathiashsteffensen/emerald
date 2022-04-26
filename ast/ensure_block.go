package ast

import (
	"bytes"
	"emerald/lexer"
)

type EnsureBlock struct {
	Token lexer.Token // The 'ensure' token
	Body  *BlockStatement
}

func (rb *EnsureBlock) TokenLiteral() string { return rb.Token.Literal }
func (rb *EnsureBlock) String() string {
	var out bytes.Buffer

	out.WriteString(rb.TokenLiteral())
	out.WriteString("\n")
	out.WriteString(rb.Body.String())

	return out.String()
}
