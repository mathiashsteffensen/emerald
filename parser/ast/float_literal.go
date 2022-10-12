package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (il *FloatLiteral) expressionNode()      {}
func (il *FloatLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *FloatLiteral) String(indents ...int) string {
	return indented(&strings.Builder{}, indents[0], il.Token.Literal).String()
}
