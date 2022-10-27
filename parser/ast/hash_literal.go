package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type HashLiteralElement struct {
	Key   Expression
	Value Expression
}

type HashLiteral struct {
	Values []*HashLiteralElement
	Token  lexer.Token // The { Token
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	if len(hl.Values) == 0 {
		indented(&out, indent, "{}")
	} else {
		indented(&out, indent, "{\n")

		for _, el := range hl.Values {
			out.WriteString(el.Key.String(indent+1) + ": " + el.Value.String(0) + ",\n")
		}

		indented(&out, indent, "}")
	}

	return out.String()
}
