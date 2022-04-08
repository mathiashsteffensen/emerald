package ast

import (
	"bytes"
	"emerald/lexer"
)

type BlockLiteral struct {
	Token      lexer.Token // The 'do' or '{' token
	Parameters []Expression
	Body       *BlockStatement
}

func (bl *BlockLiteral) expressionNode()      {}
func (bl *BlockLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BlockLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(bl.TokenLiteral())

	if len(bl.Parameters) != 0 {
		out.WriteString(" |")
		for i, parameter := range bl.Parameters {
			out.WriteString(parameter.String())
			if i != len(bl.Parameters)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString("|\n")
	}

	for _, statement := range bl.Body.Statements {
		out.WriteString("  ")
		out.WriteString(statement.String())
		out.WriteString("\n")
	}

	if bl.TokenLiteral() == "{" {
		out.WriteString("}")
	} else {
		out.WriteString("end")
	}

	return out.String()
}
