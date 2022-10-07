package ast

import (
	"bytes"
	"emerald/parser/lexer"
	"strings"
)

type MethodLiteral struct {
	Token lexer.Token // The 'def' token
	*BlockLiteral
	Name Expression
}

func (ml *MethodLiteral) expressionNode()      {}
func (ml *MethodLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MethodLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range ml.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(ml.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ml.Name.String())

	if len(params) != 0 {
		out.WriteString("(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(")")
	}

	out.WriteString("\n")
	out.WriteString(ml.Body.String())
	for _, block := range ml.RescueBlocks {
		out.WriteString(block.String())
	}
	out.WriteString("end\n")

	return out.String()
}
