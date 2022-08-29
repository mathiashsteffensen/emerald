package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestParseScopeAccessor(t *testing.T) {
	program := testParseAST(t, "MyMod::MyClass")

	expectStatementLength(t, program.Statements, 1)

	scopeAccesor, ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.ScopeAccessor)
	if !ok {
		t.Fatalf("Expression was not scope accessor got=%#v", program.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if scopeAccesor.Left.String() != "MyMod" {
		t.Errorf("expected receiver to be MyMod, got=%s", scopeAccesor.Left.String())
	}

	if scopeAccesor.Method.String() != "MyClass" {
		t.Errorf("expected receiver to be MyClass, got=%s", scopeAccesor.Method.String())
	}
}
