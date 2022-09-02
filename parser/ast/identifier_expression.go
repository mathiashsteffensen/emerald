package ast

import (
	"emerald/parser/lexer"
)

type IdentifierExpression struct {
	Token lexer.Token // the lexer.IDENT token
	Value string
}

func (ie IdentifierExpression) expressionNode()      {}
func (ie IdentifierExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IdentifierExpression) String() string       { return ie.Value }
