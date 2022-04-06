package parser

import (
	"emerald/ast"
	"testing"
)

func TestMethodLiteralExpression(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedName string
		expectedArgs []string
	}{
		{
			"with a single argument",
			`
			def method(arg)
				puts(arg)
			end
			`,
			"method",
			[]string{"arg"},
		},
		{
			"with no arguments",
			`
			def method
				puts("Hello")
			end
			`,
			"method",
			[]string{},
		},
		{
			"with multiple arguments",
			`
			def method(fmt, val)
				printf(fmt, val)
			end
			`,
			"method",
			[]string{"fmt", "val"},
		},
		{
			"one-liner",
			"def method(fmt, val); printf(fmt, val); end",
			"method",
			[]string{"fmt", "val"},
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

			literal, ok := stmt.Expression.(*ast.MethodLiteral)
			if !ok {
				t.Fatalf("exp not *ast.MethodLiteral. got=%T", stmt.Expression)
			}

			if literal.Name.TokenLiteral() != tt.expectedName {
				t.Fatalf("exp literal not %s, got=%s", tt.expectedName, literal.Name.TokenLiteral())
			}

			if len(literal.Parameters) != len(tt.expectedArgs) {
				t.Fatalf("exp %d parameters got=%d", len(tt.expectedArgs), len(literal.Parameters))
			}

			for i, parameter := range literal.Parameters {
				testIdentifier(t, parameter, tt.expectedArgs[i])
			}
		})
	}
}
