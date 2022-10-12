package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String(indents ...int) string {
	var out strings.Builder

	for _, s := range bs.Statements {
		out.WriteString(s.String(indents...))
		out.WriteString("\n")
	}

	return out.String()
}
