package ast

import (
	"emerald/lexer"
)

type BooleanExpression struct {
	Token lexer.Token
	Value bool
}

func (be *BooleanExpression) expressionNode()      {}
func (be *BooleanExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BooleanExpression) String() string       { return be.Token.Literal }
