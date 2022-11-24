package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type BlockLiteral struct {
	Token            lexer.Token // The 'do' or '{' token
	Arguments        []*IdentifierExpression
	KeywordArguments []*IdentifierExpression
	Body             *BlockStatement
	RescueBlocks     []*RescueBlock
	EnsureBlock      *EnsureBlock
}

func (bl *BlockLiteral) expressionNode()      {}
func (bl *BlockLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BlockLiteral) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]
	var newLineIndent int

	if len(indents) != 1 {
		newLineIndent = indents[1]
	} else {
		newLineIndent = indent
	}

	out.WriteString(bl.TokenLiteral())

	if len(bl.Arguments) != 0 {
		out.WriteString(" |")
		for i, parameter := range bl.Arguments {
			out.WriteString(parameter.String(0))
			if i != len(bl.Arguments)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString("|")
	}

	out.WriteString("\n")
	out.WriteString(bl.Body.String(newLineIndent + 1))

	if bl.TokenLiteral() == "{" {
		indented(&out, newLineIndent, "}")
	} else {
		for _, block := range bl.RescueBlocks {
			out.WriteString(block.String(newLineIndent))
		}

		if bl.EnsureBlock != nil {
			out.WriteString(bl.EnsureBlock.String(newLineIndent))
		}

		indented(&out, indent, "end")
	}

	return out.String()
}
