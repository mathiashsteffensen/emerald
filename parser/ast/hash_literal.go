package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type HashLiteral struct {
	Value map[Expression]Expression
	Token lexer.Token // The { Token
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	if len(hl.Value) == 0 {
		indented(&out, indent, "{}")
	} else {
		indented(&out, indent, "{\n")

		for key, value := range hl.Value {
			out.WriteString(key.String(indent+1) + ": " + value.String(0) + ",\n")
		}

		indented(&out, indent, "}")
	}

	return out.String()
}
