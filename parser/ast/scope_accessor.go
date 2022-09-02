package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type ScopeAccessor struct {
	Left  Expression
	Token lexer.Token // The :: token
	CallExpression
}

func (pa *ScopeAccessor) expressionNode()      {}
func (pa *ScopeAccessor) TokenLiteral() string { return pa.Token.Literal }
func (pa *ScopeAccessor) String() string {
	var out bytes.Buffer

	out.WriteString(pa.Left.String())
	out.WriteString(pa.TokenLiteral())
	out.WriteString(pa.Method.String())

	return out.String()
}
