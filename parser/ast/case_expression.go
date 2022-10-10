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

func (ce *CaseExpression) String() string {
	var out strings.Builder

	out.WriteString("case ")
	out.WriteString(ce.Subject.String())
	out.WriteString("\n")

	for _, clause := range ce.WhenClauses {
		out.WriteString("when ")

		matchers := []string{}

		for _, matcher := range clause.Matchers {
			matchers = append(matchers, matcher.String())
		}

		out.WriteString(strings.Join(matchers, ", "))
		out.WriteString("\n	")
		out.WriteString(clause.Consequence.String())
		out.WriteString("\n")
	}

	if ce.Alternative != nil {
		out.WriteString("else\n")
		out.WriteString(ce.Alternative.String())
		out.WriteString("\n")
	}

	out.WriteString("end\n")

	return out.String()
}
