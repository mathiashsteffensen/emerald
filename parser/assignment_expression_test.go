package parser

import (
	ast "emerald/parser/ast"
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

			testAssignmentExpression(t, stmt.Expression, tt.expectedName, tt.expectedVal)
		})
	}
}
