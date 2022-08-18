package parser

import (
	"emerald/ast"
	"testing"
)

func TestHashLiteralParsing(t *testing.T) {
	input := `{
		hello: false,
		"string": 25
	}`

	program := testParseAST(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	ok = false

	var value ast.Expression

	for key, val := range literal.Value {
		if key.String() == ":hello" {
			value = val
			ok = true
		}
	}
	if !ok {
		t.Errorf("literal.Value doesn't have expected key ':hello'. got=%q", literal.Value)
	}
	testBooleanLiteral(t, value, false)
	ok = false

	for key, val := range literal.Value {
		if key.String() == ":string" {
			value = val
			ok = true
		}
	}
	if !ok {
		t.Errorf("literal.Value doesn't have expected key ':string'. got=%q", literal.Value)
	}

	testIntegerLiteral(t, value, 25)
}
