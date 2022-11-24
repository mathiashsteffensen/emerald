package ast

import (
	"emerald/debug"
	"emerald/parser/lexer"
	"strings"
)

type InfixExpression struct {
	Token    lexer.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String(indents ...int) string {
	var out strings.Builder

	indented(&out, indents[0], "")

	if debug.IsTest {
		out.WriteRune('(')
	}

	out.WriteString(ie.Left.String(0))
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String(0, indents[0]))

	if debug.IsTest {
		out.WriteRune(')')
	}

	return out.String()
}
