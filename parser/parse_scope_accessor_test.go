package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestParseScopeAccessor(t *testing.T) {
	program := testParseAST(t, "MyMod::MyOtherMod::MyClass")

	expectStatementLength(t, program.Statements, 1)

	scopeAccessor, ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.ScopeAccessor)
	if !ok {
		t.Fatalf("Expression was not scope accessor got=%#v", program.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if left, ok := scopeAccessor.Left.(*ast.ScopeAccessor); !ok {
		t.Errorf("Left was not scope accessor got=%T", scopeAccessor.Left)
	} else {
		testIdentifier(t, left.Left, "MyMod")
		testIdentifier(t, left.Method, "MyOtherMod")
	}

	if scopeAccessor.Left.String(0) != "MyMod::MyOtherMod" {
		t.Errorf("expected receiver to be MyMod, got=%s", scopeAccessor.Left.String(0))
	}

	testIdentifier(t, scopeAccessor.Method, "MyClass")
}
