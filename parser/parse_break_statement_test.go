package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestBreakStatement(t *testing.T) {
	input := `
		break :value; break(:value)
		break if true
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 3)
	testBreakStatement(t, program.Statements[0], ":value")
	testBreakStatement(t, program.Statements[1], ":value")
	testExpressionStatement(t, program.Statements[2], func(expression *ast.IfExpression) {
		testLiteralExpression(t, expression.Condition, true)
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testBreakStatement(t, expression.Consequence.Statements[0], nil)
	})
}
