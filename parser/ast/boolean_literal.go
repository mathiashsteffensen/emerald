package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (be *BooleanLiteral) expressionNode()      {}
func (be *BooleanLiteral) TokenLiteral() string { return be.Token.Literal }
func (be *BooleanLiteral) String(indents ...int) string {
	return indented(&strings.Builder{}, indents[0], be.Token.Literal).String()
}
