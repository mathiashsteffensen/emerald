package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestArrayLiteralParsing(t *testing.T) {
	input := `[0, condition, 2]; [];

      [
        1,
      ]
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 3)

	testExpressionStatement(t, program.Statements[0], func(literal *ast.ArrayLiteral) {
		testArrayLiteral(t, literal, []any{0, "condition", 2})
	})

	testExpressionStatement(t, program.Statements[1], func(literal *ast.ArrayLiteral) {
		testArrayLiteral(t, literal, []any{})
	})

	testExpressionStatement(t, program.Statements[2], func(literal *ast.ArrayLiteral) {
		testArrayLiteral(t, literal, []any{1})
	})
}
