package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestStringLiteralExpression(t *testing.T) {
	input := `
		"hello \n\t\a\b\v\f\r\s world"
		"This is a #{template}"
		"This is a #{template} also #{boop.method} #{"nested #{template}"}"
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 3)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.StringLiteral) {
		testStringLiteral(t, expression, "hello \n\t\a\b\v\f\r  world")
	})

	testExpressionStatement(t, program.Statements[1], func(expression *ast.StringTemplate) {
		testStringLiteral(t, expression.Chain.StringLiteral, "This is a ")

		testIdentifier(t, expression.Chain.Next.Expression, "template")
	})

	testExpressionStatement(t, program.Statements[2], func(expression *ast.StringTemplate) {
		testStringLiteral(t, expression.Chain.StringLiteral, "This is a ")
		testIdentifier(t, expression.Chain.Next.Expression, "template")
		testStringLiteral(t, expression.Chain.Next.Next.StringLiteral, " also ")
		testMethodCall(t, expression.Chain.Next.Next.Next.Expression, "boop", "method", []any{}, []string{}, false)
		testStringLiteral(t, expression.Chain.Next.Next.Next.Next.StringLiteral, " ")
		testExpression(t, expression.Chain.Next.Next.Next.Next.Next.Expression, func(expression *ast.StringTemplate) {
			testStringLiteral(t, expression.Chain.StringLiteral, "nested ")
			testIdentifier(t, expression.Chain.Next.Expression, "template")
		})
	})
}
