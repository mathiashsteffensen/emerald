package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type HashLiteral struct {
	Value map[Expression]Expression
	Token lexer.Token // The { Token
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("{\n")

	for key, value := range hl.Value {
		out.WriteString("  " + key.String() + ": " + value.String() + ",\n")
	}

	out.WriteString("}")

	return out.String()
}
