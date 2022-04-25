package parser

import (
	"emerald/ast"
	"testing"
)

func TestAssignmentExpression(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedName string
		expectedVal  string
	}{
		{
			"basic assignment",
			"foobar = 5",
			"foobar",
			"5",
		},
		{
			"assignment to instance var",
			"@foobar = 5",
			"@foobar",
			"5",
		},
		{
			"assignment to global var",
			"$foobar = 5",
			"$foobar",
			"5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

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
			if ident.Name.String() != tt.expectedName {
				t.Errorf("ident.Name not %s. got=%s", tt.expectedName, ident.Name.String())
			}
			if ident.Value.String() != tt.expectedVal {
				t.Errorf("ident.Value not %s. got=%s", tt.expectedVal, ident.Value.String())
			}
		})
	}
}
