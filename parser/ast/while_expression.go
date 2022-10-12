package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type WhileExpression struct {
	Token       lexer.Token // The 'while' token
	Condition   Expression
	Consequence *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, "while ")
	out.WriteString(we.Condition.String(0))
	out.WriteString("\n")
	out.WriteString(we.Consequence.String(indent + 1))
	out.WriteString("\n")
	indented(&out, indent, "end")

	return out.String()
}
