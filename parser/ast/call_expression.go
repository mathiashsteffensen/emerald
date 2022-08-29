package ast

import (
	"bytes"
	"emerald/parser/lexer"
	"strings"
)

type CallExpression struct {
	Token     lexer.Token // The '(' token
	Method    *IdentifierExpression
	Arguments []Expression
	Block     *BlockLiteral
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Method.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
