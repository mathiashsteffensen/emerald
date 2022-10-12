package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type WhenClause struct {
	Matchers    []Expression
	Consequence *BlockStatement
}

type CaseExpression struct {
	Token       lexer.Token // The 'case' token
	Subject     Expression
	WhenClauses []*WhenClause
	Alternative *BlockStatement
}

func (ce *CaseExpression) expressionNode()      {}
func (ce *CaseExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CaseExpression) String(indents ...int) string {
	var out strings.Builder

	indent := indents[0]

	indented(&out, indent, "case ")
	out.WriteString(ce.Subject.String(0))
	out.WriteString("\n")

	for _, clause := range ce.WhenClauses {
		indented(&out, indent, "when ")

		matchers := []string{}

		for _, matcher := range clause.Matchers {
			matchers = append(matchers, matcher.String(0))
		}

		out.WriteString(strings.Join(matchers, ", "))
		out.WriteString("\n")
		out.WriteString(clause.Consequence.String(indent + 1))
		out.WriteString("\n")
	}

	if ce.Alternative != nil {
		indented(&out, indent, "else\n")
		out.WriteString(ce.Alternative.String(indent + 1))
		out.WriteString("\n")
	}

	indented(&out, indent, "end\n")

	return out.String()
}
