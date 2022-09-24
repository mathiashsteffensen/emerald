package parser

import (
	ast "emerald/parser/ast"
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
			operator:   "&&=",
			val:        false,
		},
		{
			input:      "var ||= false",
			identifier: "var",
			operator:   "||=",
			val:        false,
		},
		{
			input:      "@var &&= false",
			identifier: "@var",
			operator:   "&&=",
			val:        false,
		},
		{
			input:      "@var ||= false",
			identifier: "@var",
			operator:   "||=",
			val:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			exp := program.Statements[0].(*ast.ExpressionStatement)

			testInfixExpression(t, exp.Expression, tt.identifier, tt.operator, tt.val)
		})
	}
}
