package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestParseCaseExpression(t *testing.T) {
	input := `
		case 2
		when Integer,
			String
			true
		when 3, 4, 2.0
			5
		else
			false
		end
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.CaseExpression) {
		testLiteralExpression(t, expression.Subject, 2)

		testWhenClause(t, expression.WhenClauses[0], true, "Integer", "String")
		testWhenClause(t, expression.WhenClauses[1], 5, 3, 4, 2.0)

		expectStatementLength(t, expression.Alternative.Statements, 1)
		testExpressionStatement(t, expression.Alternative.Statements[0], func(expression ast.Expression) {
			testLiteralExpression(t, expression, false)
		})
	})
}
