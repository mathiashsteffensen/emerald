package ast

import (
	"emerald/parser/lexer"
	"strings"
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
func (ie *IfExpression) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, "if ")
	out.WriteString(ie.Condition.String(0))
	out.WriteString("\n")

	if ie.Consequence != nil {
		out.WriteString(ie.Consequence.String(indent + 1))
	} else {
		indented(&out, indent+1, "nil\n")
	}

	if ie.ElseIfs != nil {
		for _, elseIf := range ie.ElseIfs {
			indented(&out, indent, "elsif ")
			out.WriteString(elseIf.Condition.String(0))
			out.WriteString("\n")
			out.WriteString(elseIf.Consequence.String(indent + 1))
		}
	}

	if ie.Alternative != nil {
		indented(&out, indent, "else\n")
		out.WriteString(ie.Alternative.String(indent + 1))
	}

	indented(&out, indent, "end")

	return out.String()
}
