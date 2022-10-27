package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestCallExpressionParsing(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedArgs   []any
		expectedKwargs []string
		expectBlock    bool
	}{
		{
			"with parentheses & not passing a block",
			"add(1, 6, 9);",
			[]any{1, 6, 9},
			[]string{},
			false,
		},
		{
			"with parentheses and a do end block",
			`add(1, 19, 27) do
				do_stuff
			end`,
			[]any{1, 19, 27},
			[]string{},
			true,
		},
		{
			"with parentheses and a braces block",
			`add(1, 19, 27) {	do_stuff }`,
			[]any{1, 19, 27},
			[]string{},
			true,
		},
		{
			"without parentheses & not passing a block",
			"add 1, 6, 9",
			[]any{1, 6, 9},
			[]string{},
			false,
		},
		{
			"with keyword args",
			"add(base: 2, :other => 3)",
			[]any{},
			[]string{":base", ":other"},
			false,
		},
		{
			"with normal args & with keyword args",
			"add(1, 2, base: 2, :other => 3)",
			[]any{1, 2},
			[]string{":base", ":other"},
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

			testCallExpression(t, stmt.Expression, "add", tt.expectedArgs, tt.expectedKwargs, tt.expectBlock)
		})
	}
}
