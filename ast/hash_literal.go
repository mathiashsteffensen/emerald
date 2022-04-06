package ast

import (
	"bytes"
	"emerald/lexer"
)

type HashLiteral struct {
	Value map[string]Expression
	Token lexer.Token // The { Token
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("{\n")

	for key, value := range hl.Value {
		out.WriteString("  " + key + ": " + value.String() + ",\n")
	}

	out.WriteString("}")

	return out.String()
}
