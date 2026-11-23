package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestBreakStatement(t *testing.T) {
	input := `
		break :value; break(:value)
		break
		break if true
		break :value if true
		break while running
		break :value until done
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 7)
	testBreakStatement(t, program.Statements[0], ":value")
	testBreakStatement(t, program.Statements[1], ":value")
	testBreakStatement(t, program.Statements[2], nil)
	testExpressionStatement(t, program.Statements[3], func(expression *ast.IfExpression) {
		testLiteralExpression(t, expression.Condition, true)
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testBreakStatement(t, expression.Consequence.Statements[0], nil)
	})
	testExpressionStatement(t, program.Statements[4], func(expression *ast.IfExpression) {
		testLiteralExpression(t, expression.Condition, true)
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testBreakStatement(t, expression.Consequence.Statements[0], ":value")
	})
	testExpressionStatement(t, program.Statements[5], func(expression *ast.WhileExpression) {
		testIdentifier(t, expression.Condition, "running")
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testBreakStatement(t, expression.Consequence.Statements[0], nil)
	})
	testExpressionStatement(t, program.Statements[6], func(expression *ast.WhileExpression) {
		if !expression.IsUntil() {
			t.Errorf("expected an until loop")
		}

		testIdentifier(t, expression.Condition, "done")
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testBreakStatement(t, expression.Consequence.Statements[0], ":value")
	})
}

func TestBreakStatementWithChainedModifiers(t *testing.T) {
	program := testParseAST(t, "break :value if a unless b")

	expectStatementLength(t, program.Statements, 1)
	testExpressionStatement(t, program.Statements[0], func(unless *ast.IfExpression) {
		testIdentifier(t, unless.Condition, "b")
		expectStatementLength(t, unless.Alternative.Statements, 1)

		testExpressionStatement(t, unless.Alternative.Statements[0], func(ifModifier *ast.IfExpression) {
			testIdentifier(t, ifModifier.Condition, "a")
			expectStatementLength(t, ifModifier.Consequence.Statements, 1)
			testBreakStatement(t, ifModifier.Consequence.Statements[0], ":value")
		})
	})
}
