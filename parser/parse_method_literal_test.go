package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestMethodLiteralExpression(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedName         string
		expectedArgs         []string
		expectedRescueBlocks int
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
			0,
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
			0,
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
			0,
		},
		{
			"one-liner",
			"def method(fmt, val); printf(fmt, val); end",
			"method",
			[]string{"fmt", "val"},
			0,
		},
		{
			"with a rescue block",
			`
			def method
				puts("Hello")
			rescue
				puts("An error occurred")
			end
			`,
			"method",
			[]string{},
			1,
		},
		{
			"with multiple rescue blocks",
			`
			def method
				puts("Hello")
			rescue
				puts("An error occurred")
			rescue
				puts("other error occurred")
			end
			`,
			"method",
			[]string{},
			2,
		},
		{
			"with an ensure block",
			`
			def method
				puts("Hello")
			ensure
				puts(" World!")
			end
			`,
			"method",
			[]string{},
			0,
		},
		{
			"with a rescue & an ensure block",
			`
			def method
				puts("Hello")
			rescue
				puts(" Goodbye cruel")
			ensure
				puts(" World!")
			end
			`,
			"method",
			[]string{},
			1,
		},
		{
			"with multiple rescue blocks with error classes",
			`
			def method
				puts("Hello")
			rescue StandardError
				puts("An error occurred")
			rescue SystemError, NoMemoryError => e
				puts("other error occurred - " + e.inspect)
			end
			`,
			"method",
			[]string{},
			2,
		},
		{
			"assignment method",
			`
			def level=(new)
				@level = new
			end
			`,
			"level=",
			[]string{"new"},
			0,
		},
		{
			"accessing global variable in rescue block",
			`
				def suppress_argument_errors
					yield
				rescue ArgumentError
					$suppressed = $suppressed + 1
				end
			`,
			"suppress_argument_errors",
			[]string{},
			1,
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

			if len(literal.RescueBlocks) != tt.expectedRescueBlocks {
				t.Fatalf("Expected method literal to have %d rescue blocks, but got %d", len(literal.RescueBlocks), tt.expectedRescueBlocks)
			}
		})
	}
}