package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type SymbolLiteral struct {
	Token lexer.Token // The : token
	Value string
}

func (sl *SymbolLiteral) expressionNode()      {}
func (sl *SymbolLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *SymbolLiteral) String(indents ...int) string {
	return indented(&strings.Builder{}, indents[0], ":"+sl.Value).String()
}
