package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type StaticClassLiteral struct {
	Token lexer.Token // The class token
	Body  *BlockStatement
}

func (cl *StaticClassLiteral) expressionNode()      {}
func (cl *StaticClassLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *StaticClassLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("class << self")
	out.WriteString("\n")

	for _, value := range cl.Body.Statements {
		out.WriteString(value.String() + "\n")
	}

	out.WriteString("end")

	return out.String()
}
