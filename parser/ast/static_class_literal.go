package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type StaticClassLiteral struct {
	Token lexer.Token // The class token
	Body  *BlockStatement
}

func (cl *StaticClassLiteral) expressionNode()      {}
func (cl *StaticClassLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *StaticClassLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	out.WriteString(strings.Repeat("	", indent))
	out.WriteString("class << self")
	out.WriteString("\n")

	for _, value := range cl.Body.Statements {
		out.WriteString(value.String(indent+1) + "\n")
	}

	out.WriteString("end")

	return out.String()
}
