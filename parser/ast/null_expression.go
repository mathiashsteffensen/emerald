package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type NullExpression struct {
	Token lexer.Token
}

func (ne *NullExpression) expressionNode()      {}
func (ne *NullExpression) TokenLiteral() string { return ne.Token.Literal }
func (ne *NullExpression) String(indents ...int) string {
	return indented(&strings.Builder{}, indents[0], ne.Token.Literal).String()
}
