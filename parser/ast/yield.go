package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type Yield struct {
	Token            lexer.Token
	Arguments        []Expression
	KeywordArguments []*HashLiteralElement
}

func (y Yield) expressionNode()      {}
func (y Yield) TokenLiteral() string { return y.Token.Literal }
func (y Yield) String(indents ...int) string {
	var out strings.Builder

	args := []string{}
	for _, a := range y.Arguments {
		args = append(args, a.String(0))
	}

	indented(&out, indents[0], "yield(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
