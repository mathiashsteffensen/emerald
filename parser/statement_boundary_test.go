package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestNewlineEndsStatement(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"[ after call arguments", "puts foo\n[1, 2]"},
		{"( after call arguments", "puts foo\n(1)"},
		{"{ after call arguments", "puts foo\n{a: 1}"},
		{"[ after if", "if foo\n  1\nend\n[1, 2]"},
		{"[ after unless", "unless foo\n  1\nend\n[1, 2]"},
		{"[ after while", "while foo\n  1\nend\n[1, 2]"},
		{"if after call arguments", "puts foo\nif bar\n  1\nend"},
		{"if after return value", "return foo 1\nif bar\n  1\nend"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 2)
		})
	}
}

func TestModifierAppliesToWholeCall(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		method string
		args   []any
		kwargs []string
	}{
		{"after argument", `raise "x" if cond`, "raise", []any{"s:x"}, []string{}},
		{"after arguments", `raise error, "x" if cond`, "raise", []any{"error", "s:x"}, []string{}},
		{"after keyword argument", `log "x", level: 1 if cond`, "log", []any{"s:x"}, []string{":level"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)
			testExpressionStatement(t, program.Statements[0], func(modifier *ast.IfExpression) {
				testIdentifier(t, modifier.Condition, "cond")
				expectStatementLength(t, modifier.Consequence.Statements, 1)
				testExpressionStatement(t, modifier.Consequence.Statements[0], func(call ast.CallExpression) {
					testCallExpression(t, call, tt.method, tt.args, tt.kwargs, false)
				})
			})
		})
	}
}
