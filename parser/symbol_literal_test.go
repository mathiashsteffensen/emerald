package parser

import (
	"emerald/ast"
	"testing"
)

func TestSymbolLiteral(t *testing.T) {
	input := ":example"

	program := testParseAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement was not expression, got=%T", program.Statements[0])
	}

	symbol, ok := exp.Expression.(*ast.SymbolLiteral)
	if !ok {
		t.Fatalf("exp was not SymbolLiteral, got=%T", exp.Expression)
	}

	if symbol.Value != "example" {
		t.Fatalf("symbol did not have value 'example', got=%s", symbol.Value)
	}
}
