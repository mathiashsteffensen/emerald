package ast

import (
	"bytes"
	"emerald/parser/lexer"
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
func (m MethodCall) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(m.Left.String())
	out.WriteString(m.TokenLiteral())
	out.WriteString(m.Method.String())

	if len(m.Arguments) != 0 {
		out.WriteString("(")

		for i, argument := range m.Arguments {
			out.WriteString(argument.String())

			if i != len(m.Arguments)-1 {
				out.WriteString(", ")
			}
		}

		out.WriteString(")")
	}

	if m.Block != nil {
		out.WriteString(" ")
		out.WriteString(m.Block.String())
	}

	out.WriteString(")")

	return out.String()
}
