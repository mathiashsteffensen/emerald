package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type WhileExpression struct {
	Token       lexer.Token // The 'while' token
	Condition   Expression
	Consequence *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(we.Condition.String())
	out.WriteString("\n  ")
	out.WriteString(we.Consequence.String())
	out.WriteString("\nend")

	return out.String()
}
