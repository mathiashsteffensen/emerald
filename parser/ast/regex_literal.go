package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type RegexpLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *RegexpLiteral) expressionNode()      {}
func (sl *RegexpLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *RegexpLiteral) String(indents ...int) string {
	return indented(&strings.Builder{}, indents[0], `/`+sl.Value+`/`).String()
}
