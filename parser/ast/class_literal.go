package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type ClassLiteral struct {
	Token  lexer.Token // The class token
	Name   IdentifierExpression
	Parent Expression
	Body   *BlockStatement
}

func (cl *ClassLiteral) expressionNode()      {}
func (cl *ClassLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *ClassLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("class ")
	out.WriteString(cl.Name.String())
	out.WriteString("\n")

	for _, value := range cl.Body.Statements {
		out.WriteString(value.String() + "\n")
	}

	out.WriteString("end")

	return out.String()
}
