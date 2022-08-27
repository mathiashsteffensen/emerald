package parser

import (
	"emerald/ast"
	"testing"
)

func TestIdentifierExpression(t *testing.T) {
	input := "foobar; self"

	program := testParseAST(t, input)

	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	testIdentifier(t, stmt.Expression, "foobar")

	stmt, ok = program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if _, ok = stmt.Expression.(*ast.Self); !ok {
		t.Fatalf("expected ast.Self got=%T", stmt.Expression)
	}
}
