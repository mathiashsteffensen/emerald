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
		testMethodCall(t, expression, "hash", "[]", []any{":key"}, []string{}, false)
	})

	testExpressionStatement(t, program.Statements[1], func(expression *ast.MethodCall) {
		testMethodCall(t, expression, "hash", "[]=", []any{":key", ":value"}, []string{}, false)
	})

	for i, operator := range []string{"||", "&&"} {
		testExpressionStatement(t, program.Statements[i+2], func(expression *ast.AssignmentExpression) {
			testMethodCall(t, expression.Name, "cache", "[]", []any{"n"}, []string{}, false)
			value := expression.Value.(*ast.InfixExpression)
			if value.Operator != operator {
				t.Fatalf("operator = %q, want %q", value.Operator, operator)
			}
			testLiteralExpression(t, value.Right, 5)
		})
	}
}
