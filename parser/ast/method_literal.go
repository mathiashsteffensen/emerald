package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type MethodLiteral struct {
	Token lexer.Token // The 'def' token
	*BlockLiteral
	Name Expression
	// ArgumentDefaults holds the default value expressions aligned with
	// Arguments; an entry is nil when the parameter is required.
	ArgumentDefaults []Expression
}

func (ml *MethodLiteral) expressionNode()      {}
func (ml *MethodLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MethodLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	params := []string{}

	for i, p := range ml.Arguments {
		parameter := p.String(0)
		if i < len(ml.ArgumentDefaults) && ml.ArgumentDefaults[i] != nil {
			parameter += " = " + ml.ArgumentDefaults[i].String(0)
		}
		params = append(params, parameter)
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
