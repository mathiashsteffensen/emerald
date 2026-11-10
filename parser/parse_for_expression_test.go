package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestParseForExpression(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedIterable   string
		expectedStatements int
	}{
		{"on multiple lines", "for item in items\n  puts(item)\nend", "items", 1},
		{"with do", "for item in items do\n  puts(item)\nend", "items", 1},
		{"on one line", "for item in [1]; puts(item); end", "[\n\t1,\n]", 1},
		{"on one line with do", "for item in [1] do puts(item) end", "[\n\t1,\n]", 1},
		{"with an empty body", "for item in items; end", "items", 0},
		{"with a call collection and do", "for item in list(1) do\nend", "list(1)", 0},
		{"with a method call collection and do", "for item in self.items do\nend", "(self.items)", 0},
		{"with a command call collection and a newline", "for item in pair 1, 2\n  puts(item)\nend", "pair(1, 2)", 1},
		{"with a command call collection and ';'", "for item in pair 1, 2; puts(item); end", "pair(1, 2)", 1},
		{"with do and an identifier ending the body", "for item in [1] do s = item end", "[\n\t1,\n]", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			testExpressionStatement(t, program.Statements[0], func(expression *ast.ForExpression) {
				testIdentifier(t, expression.Variable, "item")

				if expression.Iterable.String(0) != tt.expectedIterable {
					t.Errorf("wrong collection, want=%s got=%s", tt.expectedIterable, expression.Iterable.String(0))
				}

				expectStatementLength(t, expression.Body.Statements, tt.expectedStatements)
			})
		})
	}
}

func TestParseForExpressionWithDoBlockInsideBracketsOfTheCollection(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"in call parentheses", "for item in wrap([3].map do |v| v end) do puts(item) end"},
		{"in a brace block", "for item in [1].map { |v| [v].each do |w| end; v } do puts(item) end"},
		{"in an array literal", "for item in [[5].map do |v| v end] do puts(item) end"},
		{"in a grouped expression", "for item in ([7].map do |v| v end) do puts(item) end"},
		{"in a hash literal", "for item in {k: [8].map do |v| v end} do puts(item) end"},
		{"in index brackets", "for item in a[[0].map do |i| i end.first] do puts(item) end"},
		{"in string interpolation", "for item in \"#{[1].map do |i| i end}\".size.to_s do puts(item) end"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			testExpressionStatement(t, program.Statements[0], func(expression *ast.ForExpression) {
				expectStatementLength(t, expression.Body.Statements, 1)
			})
		})
	}
}

func TestParseForExpressionBodyPassesDoBlocks(t *testing.T) {
	program := testParseAST(t, `
		for item in items do
			item.each do |n|
				puts(n)
			end
		end
	`)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.ForExpression) {
		testExpressionStatement(t, expression.Body.Statements[0], func(call *ast.MethodCall) {
			testMethodCall(t, call, "item", "each", []any{}, []string{}, true)
		})
	})
}

func TestParseForExpressionAfterForRestoresDoBlocks(t *testing.T) {
	program := testParseAST(t, "for a in b; end\nitems.each do |n|\nend")

	expectStatementLength(t, program.Statements, 2)

	testExpressionStatement(t, program.Statements[1], func(call *ast.MethodCall) {
		testMethodCall(t, call, "items", "each", []any{}, []string{}, true)
	})
}
