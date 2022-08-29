package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestIndexAccessor(t *testing.T) {
	input := "hash[:key]"

	program := testParseAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program did not have 1 statement, got=%d", len(program.Statements))
	}

	expr, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement was not expression statement, got=%#v", program.Statements[0])
	}

	method, ok := expr.Expression.(*ast.MethodCall)
	if !ok {
		t.Fatalf("expression was not method call, got=%#v", expr.Expression)
	}

	if method.Left.String() != "hash" {
		t.Fatalf("expected left to eq 'hash', got='%s'", method.Left)
	}

	if method.Method.String() != "[]" {
		t.Fatalf("expected left to eq '[]', got='%s'", method.Method)
	}

	if method.Arguments[0].String() != ":key" {
		t.Fatalf("expected argument to eq ':key', got='%s'", method.Arguments[0])
	}
}
