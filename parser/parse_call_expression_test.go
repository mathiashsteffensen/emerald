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

			exp, ok := stmt.Expression.(ast.CallExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
					stmt.Expression)
			}

			if !testIdentifier(t, exp.Method, "add") {
				return
			}

			if len(exp.Arguments) != len(tt.expectedArgs) {
				t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
			}

			for i, arg := range tt.expectedArgs {
				testLiteralExpression(t, exp.Arguments[i], arg)
			}

			if tt.expectBlock && exp.Block == nil {
				t.Fatalf("exp was not passed a block")
			}
		})
	}
}
