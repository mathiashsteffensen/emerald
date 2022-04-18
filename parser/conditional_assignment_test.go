package parser

import (
	"emerald/ast"
	"testing"
)

func TestConditionalAssignment(t *testing.T) {
	tests := []struct {
		input      string
		identifier string
		operator   string
		val        any
	}{
		{
			input:      "var &&= false",
			identifier: "var",
			operator:   "&&",
			val:        false,
		},
		{
			input:      "var ||= false",
			identifier: "var",
			operator:   "||",
			val:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			if len(program.Statements) != 1 {
				t.Fatalf("expected program to have %d statements, got=%d", 1, len(program.Statements))
			}

			exp := program.Statements[0].(*ast.ExpressionStatement)

			infix, ok := exp.Expression.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("expression was not *ast.InfixExpression, got=%T", exp.Expression)
			}

			if infix.Operator != tt.operator {
				t.Errorf("wrong operator want=%s, got=%s", tt.operator, infix.Operator)
			}

			testIdentifier(t, infix.Left, tt.identifier)

			assignment, ok := infix.Right.(*ast.AssignmentExpression)
			if !ok {
				t.Fatalf("expected right side of infix to be *ast.AssignmentExpression, got=%T", infix.Right)
			}

			testIdentifier(t, assignment.Name, tt.identifier)
			testLiteralExpression(t, assignment.Value, tt.val)
		})
	}
}
