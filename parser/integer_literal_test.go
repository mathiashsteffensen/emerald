package parser

import (
	"emerald/ast"
	"testing"
)

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5_000;"

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

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5000 {
		t.Errorf("literal.Value not %d. got=%d", 5000, literal.Value)
	}
	if literal.TokenLiteral() != "5_000" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5_000",
			literal.TokenLiteral())
	}
}
