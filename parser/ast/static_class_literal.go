package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type StaticClassLiteral struct {
	Token    lexer.Token // The class token
	Receiver Expression
	Body     *BlockStatement
}

func (cl *StaticClassLiteral) expressionNode()      {}
func (cl *StaticClassLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *StaticClassLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	out.WriteString(strings.Repeat("	", indent))
	out.WriteString("class << ")
	out.WriteString(cl.Receiver.String(0))
	out.WriteString("\n")

	for _, value := range cl.Body.Statements {
		out.WriteString(value.String(indent+1) + "\n")
	}

	out.WriteString("end")

	return out.String()
}
