package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestHashLiteralParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]any
	}{
		{
			input: `
				{
					hello: false,
					"string": 25
				}
			`,
			expected: map[string]any{
				":hello":  false,
				":string": 25,
			},
		},
		{
			input:    "{}",
			expected: map[string]any{},
		},
	}

	for _, tt := range tests {
		program := testParseAST(t, tt.input)

		expectStatementLength(t, program.Statements, 1)

		testExpressionStatement(t, program.Statements[0], func(expression *ast.HashLiteral) {
			testHashLiteral(t, expression, tt.expected)
		})
	}
}
