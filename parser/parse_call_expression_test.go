package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestCallExpressionParsing(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedArgs []any
		expectBlock  bool
	}{
		{
			"with parentheses & not passing a block",
			"add(1, 6, 9);",
			[]any{1, 6, 9},
			false,
		},
		{
			"with parentheses and a do end block",
			`add(1, 19, 27) do
				do_stuff
			end`,
			[]any{1, 19, 27},
			true,
		},
		{
			"with parentheses and a braces block",
			`add(1, 19, 27) {	do_stuff }`,
			[]any{1, 19, 27},
			true,
		},
		{
			"without parentheses & not passing a block",
			"add 1, 6, 9",
			[]any{1, 6, 9},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
					program.Statements[0])
			}

			testCallExpression(t, stmt.Expression, "add", tt.expectedArgs, tt.expectBlock)
		})
	}
}
