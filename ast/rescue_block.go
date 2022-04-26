package ast

import (
	"bytes"
	"emerald/lexer"
)

type RescueBlock struct {
	Token lexer.Token // The 'rescue' token
	Body  *BlockStatement
}

func (rb *RescueBlock) TokenLiteral() string { return rb.Token.Literal }
func (rb *RescueBlock) String() string {
	var out bytes.Buffer

	out.WriteString(rb.TokenLiteral())
	out.WriteString("\n")
	out.WriteString(rb.Body.String())

	return out.String()
}
