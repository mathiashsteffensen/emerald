package ast

import (
	"emerald/lexer"
)

type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (be *BooleanLiteral) expressionNode()      {}
func (be *BooleanLiteral) TokenLiteral() string { return be.Token.Literal }
func (be *BooleanLiteral) String() string       { return be.Token.Literal }
