package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type BlockLiteral struct {
	Token        lexer.Token // The 'do' or '{' token
	Parameters   []Expression
	Body         *BlockStatement
	RescueBlocks []*RescueBlock
	EnsureBlock  *EnsureBlock
}

func (bl *BlockLiteral) expressionNode()      {}
func (bl *BlockLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BlockLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	out.WriteString(bl.TokenLiteral())

	if len(bl.Parameters) != 0 {
		out.WriteString(" |")
		for i, parameter := range bl.Parameters {
			out.WriteString(parameter.String(0))
			if i != len(bl.Parameters)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString("|")
	}

	out.WriteString("\n")
	out.WriteString(bl.Body.String(indent + 1))

	if bl.TokenLiteral() == "{" {
		indented(&out, indent, "}")
	} else {
		for _, block := range bl.RescueBlocks {
			out.WriteString(block.String(indent))
		}

		if bl.EnsureBlock != nil {
			out.WriteString(bl.EnsureBlock.String(indent))
		}

		indented(&out, indent, "end")
	}

	return out.String()
}
