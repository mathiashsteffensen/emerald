package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestParseRegexpLiteral(t *testing.T) {
	input := "/r/"

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.RegexpLiteral) {
		if expression.Value != "r" {
			t.Fatalf("Expected regexp value to be %q but got %q", "r", expression.Value)
		}
	})
}
