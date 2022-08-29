package ast

import "emerald/lexer"

type RegexpLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *RegexpLiteral) expressionNode()      {}
func (sl *RegexpLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *RegexpLiteral) String() string       { return `/` + sl.Value + `/` }
