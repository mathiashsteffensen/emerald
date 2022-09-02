package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestParseYield(t *testing.T) {
	program := testParseAST(t, "yield(1, 2)")

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(expression ast.Yield) {
		testIntegerLiteral(t, expression.Arguments[0], 1)
		testIntegerLiteral(t, expression.Arguments[1], 2)
	})
}
