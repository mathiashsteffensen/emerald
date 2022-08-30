package ast

import (
	"emerald/parser/lexer"
)

type SymbolLiteral struct {
	Token lexer.Token // The : token
	Value string
}

func (sl *SymbolLiteral) expressionNode()      {}
func (sl *SymbolLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *SymbolLiteral) String() string       { return ":" + sl.Value }
