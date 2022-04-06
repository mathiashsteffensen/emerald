package parser

import (
	"emerald/ast"
	"testing"
)

func TestNullExpression(t *testing.T) {
	input := "nil"

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

	_, ok = stmt.Expression.(*ast.NullExpression)
	if !ok {
		t.Fatalf("exp not *ast.NullExpression. got=%T", stmt.Expression)
	}
}
