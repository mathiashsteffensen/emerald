package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type CallExpression struct {
	Token            lexer.Token // The '(' token
	Method           IdentifierExpression
	Arguments        []Expression
	KeywordArguments []*HashLiteralElement
	Block            *BlockLiteral
}

func (ce CallExpression) expressionNode()      {}
func (ce CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce CallExpression) String(indents ...int) string {
	var out strings.Builder

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String(0))
	}
	for _, el := range ce.KeywordArguments {
		args = append(args, strings.Join([]string{el.Key.String(0), el.Value.String(0)}, ": "))
	}

	out.WriteString(ce.Method.String(indents...))
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	if ce.Block != nil {
		var blockIndent int
		if len(indents) != 1 {
			blockIndent = indents[1]
		} else {
			blockIndent = indents[0]
		}

		out.WriteString(" ")
		out.WriteString(ce.Block.String(blockIndent))
	}

	return out.String()
}
