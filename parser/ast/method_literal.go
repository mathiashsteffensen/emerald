package ast

import (
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
func (ml *MethodLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	params := []string{}

	for _, p := range ml.Arguments {
		params = append(params, p.String(0))
	}

	indented(&out, indent, "def ")
	out.WriteString(ml.Name.String(0))

	if len(params) != 0 {
		out.WriteString("(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(")")
	}

	out.WriteString("\n")
	out.WriteString(ml.Body.String(indent + 1))
	for _, block := range ml.RescueBlocks {
		out.WriteString(block.String(indent))
	}

	indented(&out, indent, "end")

	return out.String()
}
