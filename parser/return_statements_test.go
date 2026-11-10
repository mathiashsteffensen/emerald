package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue any
	}{
		{"returning a number", "return 5;", 5},
		{"returning a boolean", "return true;", true},
		{"returning an condition", "return foobar;", "foobar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statements. got=%d",
					len(program.Statements))
			}

			stmt := program.Statements[0]
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok {
				t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
			}
			if returnStmt.TokenLiteral() != "return" {
				t.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
			}

			testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue)
		})
	}
}

func TestReturnStatementWithModifier(t *testing.T) {
	program := testParseAST(t, "return foo 1 if cond")

	expectStatementLength(t, program.Statements, 1)
	testExpressionStatement(t, program.Statements[0], func(modifier *ast.IfExpression) {
		testIdentifier(t, modifier.Condition, "cond")
		expectStatementLength(t, modifier.Consequence.Statements, 1)

		returnStmt, ok := modifier.Consequence.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("modifier statement is not *ast.ReturnStatement. got=%T", modifier.Consequence.Statements[0])
		}

		testCallExpression(t, returnStmt.ReturnValue, "foo", []any{1}, []string{}, false)
	})
}
