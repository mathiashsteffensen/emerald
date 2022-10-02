package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type BreakStatement struct {
	Token lexer.Token // break token
	Value Expression
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string {
	var out bytes.Buffer

	out.WriteString(bs.TokenLiteral())

	if bs.Value != nil {
		out.WriteString(" ")
		out.WriteString(bs.Value.String())
	}

	return out.String()
}
