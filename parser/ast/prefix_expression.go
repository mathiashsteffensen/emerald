package ast

import (
	"emerald/debug"
	"emerald/parser/lexer"
	"strings"
)

type PrefixExpression struct {
	Token    lexer.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String(indents ...int) string {
	var out strings.Builder

	indented(&out, indents[0], "")

	if debug.IsTest {
		out.WriteRune('(')
	}

	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String(0))

	if debug.IsTest {
		out.WriteRune(')')
	}

	return out.String()
}
