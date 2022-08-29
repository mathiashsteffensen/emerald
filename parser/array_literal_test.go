package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestArrayLiteralParsing(t *testing.T) {
	input := `[0, identifier, 2]`

	expected := []any{0, "identifier", 2}

	program := testParseAST(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(literal.Value) != len(expected) {
		t.Fatalf("exp does not %d values got=%d", len(expected), len(literal.Value))
	}

	for i, exp := range expected {
		testLiteralExpression(t, literal.Value[i], exp)
	}
}
