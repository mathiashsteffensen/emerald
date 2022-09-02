package ast

import (
	"bytes"
	"emerald/parser/lexer"
	"strings"
)

type Yield struct {
	Token     lexer.Token
	Arguments []Expression
}

func (y Yield) expressionNode()      {}
func (y Yield) TokenLiteral() string { return y.Token.Literal }
func (y Yield) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range y.Arguments {
		args = append(args, a.String())
	}

	out.WriteString("yield(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
