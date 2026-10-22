package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestNextStatement(t *testing.T) {
	input := `
		next; next :value
		next if true
		next :value unless true
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 4)
	testNextStatement(t, program.Statements[0], nil)
	testNextStatement(t, program.Statements[1], ":value")
	testExpressionStatement(t, program.Statements[2], func(expression *ast.IfExpression) {
		testLiteralExpression(t, expression.Condition, true)
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testNextStatement(t, expression.Consequence.Statements[0], nil)
	})
	testExpressionStatement(t, program.Statements[3], func(expression *ast.IfExpression) {
		expectStatementLength(t, expression.Alternative.Statements, 1)
		testNextStatement(t, expression.Alternative.Statements[0], ":value")
	})
}
