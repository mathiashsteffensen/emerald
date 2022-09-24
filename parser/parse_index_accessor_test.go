package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestIndexAccessor(t *testing.T) {
	input := `hash[:key]; cache[n] ||= 5`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 2)
	testExpressionStatement(t, program.Statements[0], func(expression *ast.MethodCall) {
		testMethodCall(t, expression, "hash", "[]", []any{":key"}, false)
	})
	testExpressionStatement(t, program.Statements[1], func(expression *ast.InfixExpression) {
		testMethodCall(t, expression.Left, "cache", "[]", []any{"n"}, false)
		if expression.Operator != "||=" {
			t.Fatalf("Expected operator to be ||= but got %s", expression.Operator)
		}
		testIntegerLiteral(t, expression.Right, 5)
	})
}
