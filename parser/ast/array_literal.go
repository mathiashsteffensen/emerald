package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type ArrayLiteral struct {
	Value []Expression
	Token lexer.Token // The [ token
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	if len(al.Value) == 0 {
		return indented(&out, indent, "[]").String()
	}

	indented(&out, indent, "[\n")

	var newlineIndent int
	if len(indents) != 1 {
		newlineIndent = indents[1]
	} else {
		newlineIndent = indents[0]
	}

	for _, value := range al.Value {
		out.WriteString(value.String(newlineIndent+1) + ",\n")
	}

	indented(&out, newlineIndent, "]")

	return out.String()
}
