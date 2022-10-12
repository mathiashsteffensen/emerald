package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type MethodCall struct {
	Left  Expression
	Token lexer.Token // The . token
	CallExpression
}

func (m MethodCall) Dup() *MethodCall {
	return &MethodCall{
		Left:           m.Left,
		Token:          m.Token,
		CallExpression: m.CallExpression,
	}
}
func (m MethodCall) TokenLiteral() string { return m.Token.Literal }
func (m MethodCall) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, "(")
	out.WriteString(m.Left.String(0))
	out.WriteString(".")
	out.WriteString(m.Method.String(0))

	if len(m.Arguments) != 0 {
		out.WriteString("(")

		for i, argument := range m.Arguments {
			out.WriteString(argument.String(0))

			if i != len(m.Arguments)-1 {
				out.WriteString(", ")
			}
		}

		out.WriteString(")")
	}

	if m.Block != nil {
		var blockIndent int
		if len(indents) != 1 {
			blockIndent = indents[1]
		} else {
			blockIndent = indents[0]
		}

		out.WriteString(" ")
		out.WriteString(m.Block.String(blockIndent))
	}

	out.WriteString(")")

	return out.String()
}
