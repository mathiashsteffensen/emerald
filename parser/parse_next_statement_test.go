package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestNextStatement(t *testing.T) {
	input := `
		next; next :value
		next if true
		next :value unless true
		next while running
		next :value until done
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 6)
	testNextStatement(t, program.Statements[0], nil)
	testNextStatement(t, program.Statements[1], ":value")
	testExpressionStatement(t, program.Statements[2], func(expression *ast.IfExpression) {
		testLiteralExpression(t, expression.Condition, true)
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testNextStatement(t, expression.Consequence.Statements[0], nil)
	})
	testExpressionStatement(t, program.Statements[3], func(expression *ast.IfExpression) {
		expectStatementLength(t, expression.Alternative.Statements, 1)
		testNextStatement(t, expression.Alternative.Statements[0], ":value")
	})
	testExpressionStatement(t, program.Statements[4], func(expression *ast.WhileExpression) {
		testIdentifier(t, expression.Condition, "running")
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testNextStatement(t, expression.Consequence.Statements[0], nil)
	})
	testExpressionStatement(t, program.Statements[5], func(expression *ast.WhileExpression) {
		if !expression.IsUntil() {
			t.Errorf("expected an until loop")
		}

		testIdentifier(t, expression.Condition, "done")
		expectStatementLength(t, expression.Consequence.Statements, 1)
		testNextStatement(t, expression.Consequence.Statements[0], ":value")
	})
}

func TestNextStatementWithChainedModifiers(t *testing.T) {
	program := testParseAST(t, "next if c while d")

	expectStatementLength(t, program.Statements, 1)
	testExpressionStatement(t, program.Statements[0], func(loop *ast.WhileExpression) {
		testIdentifier(t, loop.Condition, "d")
		expectStatementLength(t, loop.Consequence.Statements, 1)

		testExpressionStatement(t, loop.Consequence.Statements[0], func(ifModifier *ast.IfExpression) {
			testIdentifier(t, ifModifier.Condition, "c")
			expectStatementLength(t, ifModifier.Consequence.Statements, 1)
			testNextStatement(t, ifModifier.Consequence.Statements[0], nil)
		})
	})
}
