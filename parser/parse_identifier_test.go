package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestParseIdentifier(t *testing.T) {
	input := "foobar; self; method { 2 }"

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 3)

	testExpressionStatement(t, program.Statements[0], func(expression ast.IdentifierExpression) {
		testIdentifier(t, expression, "foobar")
	})

	testExpressionStatement(t, program.Statements[1], func(_ *ast.Self) {})
}
