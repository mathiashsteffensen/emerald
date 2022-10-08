package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type ElseIf struct {
	Condition   Expression
	Consequence *BlockStatement
}

type IfExpression struct {
	Token       lexer.Token // The 'if' or 'unless' token
	Condition   Expression
	Consequence *BlockStatement
	ElseIfs     []ElseIf
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString("\n	")

	if ie.Consequence != nil {
		out.WriteString(ie.Consequence.String())
	} else {
		out.WriteString("	nil")
	}

	if ie.ElseIfs != nil {
		out.WriteString("\n")
		for _, elseIf := range ie.ElseIfs {
			out.WriteString("elsif ")
			out.WriteString(elseIf.Condition.String())
			out.WriteString("\n	")
			out.WriteString(elseIf.Consequence.String())
		}
	}

	if ie.Alternative != nil {
		out.WriteString("\nelse\n  ")
		out.WriteString(ie.Alternative.String())
	}

	out.WriteString("\nend")

	return out.String()
}
