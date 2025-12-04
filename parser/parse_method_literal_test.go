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
		expectedKwargs       []string
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
			[]string{},
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
			[]string{},
			0,
		},
		{
			"with no arguments with parens",
			`
			def method()
				puts("Hello")
			end
			`,
			"method",
			[]string{},
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
			[]string{},
			0,
		},
		{
			"one-liner",
			"def method(fmt, val); printf(fmt, val); end",
			"method",
			[]string{"fmt", "val"},
			[]string{},
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
			[]string{},
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
			[]string{},
			1,
		},
		{
			name:                 "with keyword arguments",
			input:                "def method(option:, other:); option; end",
			expectedName:         "method",
			expectedArgs:         []string{},
			expectedKwargs:       []string{"option", "other"},
			expectedRescueBlocks: 0,
		},
		{
			name:                 "with arguments & keyword arguments",
			input:                "def method(name, option:); name + option; end",
			expectedName:         "method",
			expectedArgs:         []string{"name"},
			expectedKwargs:       []string{"option"},
			expectedRescueBlocks: 0,
		},
		{
			name:                 "with trailing comma in parameters",
			input:                "def method(a,); end",
			expectedName:         "method",
			expectedArgs:         []string{"a"},
			expectedKwargs:       []string{},
			expectedRescueBlocks: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			testExpressionStatement(t, program.Statements[0], func(literal *ast.MethodLiteral) {
				testLiteralExpression(t, literal.Name, tt.expectedName)

				if len(literal.Arguments) != len(tt.expectedArgs) {
					t.Fatalf("exp %d parameters got=%d", len(tt.expectedArgs), len(literal.Arguments))
				}

				for i, parameter := range literal.Arguments {
					testIdentifier(t, *parameter, tt.expectedArgs[i])
				}

				for i, argument := range literal.KeywordArguments {
					testIdentifier(t, *argument, tt.expectedKwargs[i])
				}

				if len(literal.RescueBlocks) != tt.expectedRescueBlocks {
					t.Fatalf("Expected method literal to have %d rescue blocks, but got %d", len(literal.RescueBlocks), tt.expectedRescueBlocks)
				}
			})
		})
	}

}
