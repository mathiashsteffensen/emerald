package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type ScopeAccessor struct {
	Left  Expression
	Token lexer.Token // The :: token
	CallExpression
}

func (pa *ScopeAccessor) expressionNode()      {}
func (pa *ScopeAccessor) TokenLiteral() string { return pa.Token.Literal }
func (pa *ScopeAccessor) String(indents ...int) string {
	var out strings.Builder

	out.WriteString(pa.Left.String(0))
	out.WriteString(pa.TokenLiteral())
	out.WriteString(pa.Method.String(0))

	return out.String()
}
