package ast

import (
	"emerald/parser/lexer"
	"strings"
)

// WhileExpression is a `while` or an `until` loop, told apart by its token
type WhileExpression struct {
	Token       lexer.Token // The 'while' or 'until' token
	Condition   Expression
	Consequence *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }

// IsUntil reports whether the loop runs while its condition is falsy
func (we *WhileExpression) IsUntil() bool { return we.Token.Type == lexer.UNTIL }

func (we *WhileExpression) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, we.Token.Literal+" ")
	out.WriteString(we.Condition.String(0))
	out.WriteString("\n")
	out.WriteString(we.Consequence.String(indent + 1))
	indented(&out, indent, "end")

	return out.String()
}
