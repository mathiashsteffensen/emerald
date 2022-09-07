package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestParseIdentifier(t *testing.T) {
	input := `foobar; self; method { 2 }
		method do
			2
		end
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 4)

	testExpressionStatement(t, program.Statements[0], func(expression ast.IdentifierExpression) {
		testIdentifier(t, expression, "foobar")
	})

	testExpressionStatement(t, program.Statements[1], func(_ *ast.Self) {})

	testExpressionStatement(t, program.Statements[2], func(expression ast.CallExpression) {
		testIdentifier(t, expression.Method, "method")
		expectStatementLength(t, expression.Block.Body.Statements, 1)
	})
}
