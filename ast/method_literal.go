package ast

import (
	"bytes"
	"emerald/lexer"
	"strings"
)

type MethodLiteral struct {
	Token      lexer.Token // The 'def' token
	Parameters []Expression
	Body       *BlockStatement
	Name       Expression
}

func (fl *MethodLiteral) expressionNode()      {}
func (fl *MethodLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *MethodLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}
