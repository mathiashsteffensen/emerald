package ast

import (
	"emerald/debug"
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
	var newLineIndent int

	if len(indents) != 1 {
		newLineIndent = indents[1]
	} else {
		newLineIndent = indent
	}

	indented(&out, indent, "")

	if debug.IsTest {
		out.WriteRune('(')
	}

	out.WriteString(m.Left.String(0, newLineIndent))
	out.WriteString(".")
	out.WriteString(m.Method.String(0, newLineIndent))

	if len(m.Arguments) != 0 {
		out.WriteString("(")

		args := []string{}
		for _, a := range m.Arguments {
			args = append(args, a.String(0))
		}
		for _, el := range m.KeywordArguments {
			args = append(args, el.Key.String(0)+": "+el.Value.String(0))
		}

		out.WriteString(strings.Join(args, ", "))

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
		out.WriteString(m.Block.String(blockIndent, newLineIndent))
	}

	if debug.IsTest {
		out.WriteRune(')')
	}

	return out.String()
}
