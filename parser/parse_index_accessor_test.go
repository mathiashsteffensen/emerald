package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestIndexAccessor(t *testing.T) {
	input := `
		hash[:key]
		hash[:key] = :value
		cache[n] ||= 5
		cache[n] &&= 5
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 4)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.MethodCall) {
		testMethodCall(t, expression, "hash", "[]", []any{":key"}, false)
	})

	testExpressionStatement(t, program.Statements[1], func(expression *ast.MethodCall) {
		testMethodCall(t, expression, "hash", "[]=", []any{":key", ":value"}, false)
	})

	testExpressionStatement(t, program.Statements[2], func(expression *ast.IfExpression) {
		testMethodCall(t, expression.Condition, "cache", "[]", []any{"n"}, false)

		testExpressionStatement(t, expression.Consequence.Statements[0], func(expression *ast.MethodCall) {
			testMethodCall(t, expression, "cache", "[]", []any{"n"}, false)

		})

		testExpressionStatement(t, expression.Alternative.Statements[0], func(expression *ast.MethodCall) {
			testMethodCall(t, expression, "cache", "[]=", []any{"n", 5}, false)
		})
	})

	testExpressionStatement(t, program.Statements[3], func(expression *ast.IfExpression) {
		testMethodCall(t, expression.Condition, "cache", "[]", []any{"n"}, false)

		testExpressionStatement(t, expression.Consequence.Statements[0], func(expression *ast.MethodCall) {
			testMethodCall(t, expression, "cache", "[]=", []any{"n", 5}, false)
		})

		testExpressionStatement(t, expression.Alternative.Statements[0], func(expression *ast.MethodCall) {
			testMethodCall(t, expression, "cache", "[]", []any{"n"}, false)
		})
	})
}
