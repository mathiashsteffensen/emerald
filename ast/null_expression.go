package ast

import "emerald/lexer"

type NullExpression struct {
	Token lexer.Token
}

func (ne *NullExpression) expressionNode()      {}
func (ne *NullExpression) TokenLiteral() string { return ne.Token.Literal }
func (ne *NullExpression) String() string       { return ne.Token.Literal }
