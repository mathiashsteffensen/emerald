package parser

import (
	"emerald/ast"
	"testing"
)

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	program := testParseAST(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}
