package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String(indents ...int) string {
	return indented(&strings.Builder{}, indents[0], `"`+sl.Value+`"`).String()
}
