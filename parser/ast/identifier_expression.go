package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type IdentifierExpression struct {
	Token lexer.Token // the lexer.IDENT token
	Value string
}

func (ie IdentifierExpression) expressionNode()      {}
func (ie IdentifierExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IdentifierExpression) String(indents ...int) string {
	return indented(&strings.Builder{}, indents[0], ie.Value).String()
}
