package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestStringLiteralExpression(t *testing.T) {
	input := `
		"hello world"
		"This is a #{template}"
		"This is a #{template} also #{boop.method}"
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 3)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.StringLiteral) {
		testStringLiteral(t, expression, "hello world")
	})

	testExpressionStatement(t, program.Statements[1], func(expression *ast.StringTemplate) {
		testStringLiteral(t, expression.Chain.StringLiteral, "This is a ")

		testIdentifier(t, expression.Chain.Next.Expression, "template")
	})

	testExpressionStatement(t, program.Statements[2], func(expression *ast.StringTemplate) {
		testStringLiteral(t, expression.Chain.StringLiteral, "This is a ")
		testIdentifier(t, expression.Chain.Next.Expression, "template")
		testStringLiteral(t, expression.Chain.Next.Next.StringLiteral, " also ")
		testMethodCall(t, expression.Chain.Next.Next.Next.Expression, "boop", "method", []any{}, false)
	})
}
