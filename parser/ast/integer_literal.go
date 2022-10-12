package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String(indents ...int) string {
	return indented(&strings.Builder{}, indents[0], il.Token.Literal).String()
}
