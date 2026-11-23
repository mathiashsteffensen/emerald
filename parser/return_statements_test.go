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

func TestReturnStatementWithLoopModifier(t *testing.T) {
	program := testParseAST(t, "return 1 until done; return while running")

	expectStatementLength(t, program.Statements, 2)
	testExpressionStatement(t, program.Statements[0], func(modifier *ast.WhileExpression) {
		if !modifier.IsUntil() {
			t.Errorf("expected an until loop")
		}

		testIdentifier(t, modifier.Condition, "done")
		expectStatementLength(t, modifier.Consequence.Statements, 1)

		returnStmt, ok := modifier.Consequence.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("modifier statement is not *ast.ReturnStatement. got=%T", modifier.Consequence.Statements[0])
		}

		testLiteralExpression(t, returnStmt.ReturnValue, 1)
	})
	testExpressionStatement(t, program.Statements[1], func(modifier *ast.WhileExpression) {
		if modifier.IsUntil() {
			t.Errorf("expected a while loop")
		}

		testIdentifier(t, modifier.Condition, "running")
		expectStatementLength(t, modifier.Consequence.Statements, 1)

		if _, ok := modifier.Consequence.Statements[0].(*ast.ReturnStatement); !ok {
			t.Fatalf("modifier statement is not *ast.ReturnStatement. got=%T", modifier.Consequence.Statements[0])
		}
	})
}

func TestReturnStatementWithChainedModifiers(t *testing.T) {
	program := testParseAST(t, "return 1 if a unless b")

	expectStatementLength(t, program.Statements, 1)

	if program.Statements[0].TokenLiteral() != "return" {
		t.Errorf("wrong token literal, got=%q", program.Statements[0].TokenLiteral())
	}

	testExpressionStatement(t, program.Statements[0], func(unless *ast.IfExpression) {
		testIdentifier(t, unless.Condition, "b")
		expectStatementLength(t, unless.Alternative.Statements, 1)

		testExpressionStatement(t, unless.Alternative.Statements[0], func(ifModifier *ast.IfExpression) {
			testIdentifier(t, ifModifier.Condition, "a")
			expectStatementLength(t, ifModifier.Consequence.Statements, 1)

			returnStmt, ok := ifModifier.Consequence.Statements[0].(*ast.ReturnStatement)
			if !ok {
				t.Fatalf("modifier statement is not *ast.ReturnStatement. got=%T", ifModifier.Consequence.Statements[0])
			}

			testLiteralExpression(t, returnStmt.ReturnValue, 1)
		})
	})
}
