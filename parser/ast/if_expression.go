package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type IfExpression struct {
	Token       lexer.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString("\n  ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("\nelse\n  ")
		out.WriteString(ie.Alternative.String())
	}

	out.WriteString("\nend")

	return out.String()
}
