package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestConditionalAssignment(t *testing.T) {
	tests := []struct {
		input       string
		condition   string
		consequence string
		alternative string
	}{
		{
			input:       "var &&= false",
			condition:   "var",
			consequence: "(var = false)",
			alternative: "var",
		},
		{
			input:       "var ||= false",
			condition:   "var",
			consequence: "var",
			alternative: "(var = false)",
		},
		{
			input:       "@var &&= false",
			condition:   "@var",
			consequence: "(@var = false)",
			alternative: "@var",
		},
		{
			input:       "@var ||= false",
			condition:   "@var",
			consequence: "@var",
			alternative: "(@var = false)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			testExpressionStatement(t, program.Statements[0], func(expression *ast.IfExpression) {
				testIfElseExpression(t, expression, tt.condition, tt.consequence, tt.alternative)
			})
		})
	}
}
