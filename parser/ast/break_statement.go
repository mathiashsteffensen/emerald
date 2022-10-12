package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type BreakStatement struct {
	Token lexer.Token // break token
	Value Expression
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String(indents ...int) string {
	var out strings.Builder

	indented(&out, indents[0], bs.TokenLiteral())

	if bs.Value != nil {
		out.WriteString(" ")
		out.WriteString(bs.Value.String(0))
	}

	return out.String()
}
