package parser

import (
	"emerald/ast"
	"testing"
)

func TestBlockLiteral(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedArgs []string
	}{
		{
			"with '{' delimiter",
			"[0].map { |n| n+1 }",
			[]string{"n"},
		},
		{
			"with 'do' & 'end' delimiters",
			`
			[0].map do |n|
				n + 1
			end`,
			[]string{"n"},
		},
		{
			"with 'do' & 'end' delimiters when block is passed to call expression",
			`
			[0].map do |n|
				n + 1
			end`,
			[]string{"n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
			}

			stmt := program.Statements[0]

			expStmt, ok := stmt.(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("statement is not ExpressionStatement, got=%T", stmt)
			}

			methodCall, ok := expStmt.Expression.(*ast.MethodCall)
			if !ok {
				t.Fatalf("expression was not MethodCall, got=%T", expStmt.Expression)
			}

			if methodCall.Block == nil {
				t.Fatalf("method call was not passed a block")
			}

			for i, arg := range tt.expectedArgs {
				testIdentifier(t, methodCall.Block.Parameters[i], arg)
			}
		})
	}
}
