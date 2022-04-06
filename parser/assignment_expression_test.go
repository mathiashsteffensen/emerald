package parser

import (
	"emerald/ast"
	"testing"
)

func TestAssignmentExpression(t *testing.T) {
	input := "foobar = 5"

	program := testParseAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("exp not *ast.AssignmentExpression. got=%T", stmt.Expression)
	}
	if ident.Name.String() != "foobar" {
		t.Errorf("ident.Name not %s. got=%s", "foobar", ident.Name.String())
	}
	if ident.Value.String() != "5" {
		t.Errorf("ident.Value not %s. got=%s", "5", ident.Value.String())
	}
}
