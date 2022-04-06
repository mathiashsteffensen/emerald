package parser

import (
	"emerald/ast"
	"testing"
)

func TestMethodCallParsing(t *testing.T) {
	input := "1.add(2, 3, 4 + 5);"

	program := testParseAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.MethodCall)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.MethodCall. got=%T", stmt.Expression)
	}

	if !testLiteralExpression(t, exp.Left, 1) {
		return
	}

	if !testIdentifier(t, exp.Method, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 2)
	testLiteralExpression(t, exp.Arguments[1], 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}
