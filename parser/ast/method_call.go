package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type MethodCall struct {
	Left  Expression
	Token lexer.Token // The . token
	*CallExpression
}

func (pa *MethodCall) expressionNode()      {}
func (pa *MethodCall) TokenLiteral() string { return pa.Token.Literal }
func (pa *MethodCall) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pa.Left.String())
	out.WriteString(pa.TokenLiteral())
	out.WriteString(pa.Method.String())

	if len(pa.Arguments) != 0 {
		out.WriteString("(")

		for i, argument := range pa.Arguments {
			out.WriteString(argument.String())

			if i != len(pa.Arguments)-1 {
				out.WriteString(", ")
			}
		}

		out.WriteString(")")
	}

	if pa.Block != nil {
		out.WriteString(" ")
		out.WriteString(pa.Block.String())
	}

	out.WriteString(")")

	return out.String()
}
