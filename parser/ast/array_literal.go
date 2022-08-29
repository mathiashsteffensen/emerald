package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type ArrayLiteral struct {
	Value []Expression
	Token lexer.Token // The [ token
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[\n")

	for _, value := range al.Value {
		out.WriteString(value.String() + ",\n")
	}

	out.WriteString("]")

	return out.String()
}
