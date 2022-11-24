package ast

import (
	"emerald/debug"
	"emerald/parser/lexer"
	"strings"
)

type AssignmentExpression struct {
	Token lexer.Token // the lexer.IDENT token
	Name  Expression
	Value Expression
}

func (ls *AssignmentExpression) expressionNode()      {}
func (ls *AssignmentExpression) TokenLiteral() string { return ls.Token.Literal }
func (ls *AssignmentExpression) String(indents ...int) string {
	var out strings.Builder

	indented(&out, indents[0], "")

	if debug.IsTest {
		out.WriteRune('(')
	}

	out.WriteString(ls.Name.String(0))
	out.WriteString(" = ")
	out.WriteString(ls.Value.String(0, indents[0], indents[0]))

	if debug.IsTest {
		out.WriteRune(')')
	}

	return out.String()
}
