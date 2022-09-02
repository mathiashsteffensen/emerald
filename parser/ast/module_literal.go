package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type ModuleLiteral struct {
	Token lexer.Token // The module token
	Body  *BlockStatement
	Name  IdentifierExpression
}

func (cl *ModuleLiteral) expressionNode()      {}
func (cl *ModuleLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *ModuleLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("module ")
	out.WriteString(cl.Name.String())
	out.WriteString("\n")

	for _, value := range cl.Body.Statements {
		out.WriteString(value.String() + "\n")
	}

	out.WriteString("end")

	return out.String()
}
