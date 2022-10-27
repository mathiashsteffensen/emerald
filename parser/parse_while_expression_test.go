package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestParseWhileExpression(t *testing.T) {
	input := `
		items = [1, 4, 9]
		while item = items.pop
			puts(item)
		end
		puts("Hello") while items.pop
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 3)

	testExpressionStatement(t, program.Statements[1], func(expression *ast.WhileExpression) {
		testAssignmentExpression(t, expression.Condition, "item", "(items.pop)")
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testExpressionStatement(t, expression.Consequence.Statements[0], func(expression ast.CallExpression) {
			testIdentifier(t, expression.Method, "puts")
			testIdentifier(t, expression.Arguments[0], "item")
		})
	})

	testExpressionStatement(t, program.Statements[2], func(expression *ast.WhileExpression) {
		testMethodCall(t, expression.Condition, "items", "pop", []any{}, []string{}, false)
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testExpressionStatement(t, expression.Consequence.Statements[0], func(expression ast.CallExpression) {
			testCallExpression(t, expression, "puts", []any{"s:Hello"}, []string{}, false)
		})
	})
}
