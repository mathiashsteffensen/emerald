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
			"with a single argument and no parens",
			`
			def greet name
				puts(name)
			end
			`,
			"greet",
			[]string{"name"},
			[]string{},
			0,
		},
		{
			"with multiple arguments and no parens",
			`
			def add left, right
				left + right
			end
			`,
			"add",
			[]string{"left", "right"},
			[]string{},
			0,
		},
		{
			"with multiple arguments across lines and no parens",
			`
			def add left,
				right
				left + right
			end
			`,
			"add",
			[]string{"left", "right"},
			[]string{},
			0,
		},
		{
			"with a positional argument and a keyword argument across lines and no parens",
			`
			def add left,
				right:
				left
			end
			`,
			"add",
			[]string{"left"},
			[]string{"right"},
			0,
		},
		{
			"with keyword arguments across lines and no parens",
			`
			def add left,
				right:,
				other:
				left
			end
			`,
			"add",
			[]string{"left"},
			[]string{"right", "other"},
			0,
		},
		{
			"with keyword arguments and no parens",
			"def method name, option:; name; end",
			"method",
			[]string{"name"},
			[]string{"option"},
			0,
		},
		{
			"with only keyword arguments and no parens",
			"def method option:; option; end",
			"method",
			[]string{},
			[]string{"option"},
			0,
		},
		{
			"with multiple keyword arguments and no parens",
			"def method first:, second:; first; end",
			"method",
			[]string{},
			[]string{"first", "second"},
			0,
		},
		{
			"with keywords as keyword argument names",
			"def method(in:, for:); end",
			"method",
			[]string{},
			[]string{"in", "for"},
			0,
		},
		{
			"with keywords as keyword argument names and no parens",
			"def method a, in:, for:; end",
			"method",
			[]string{"a"},
			[]string{"in", "for"},
			0,
		},
		{
			"with multiple arguments and no parens on one line",
			"def add left, right; left + right; end",
			"add",
			[]string{"left", "right"},
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

				if len(literal.KeywordArguments) != len(tt.expectedKwargs) {
					t.Fatalf("exp %d keyword parameters got=%d", len(tt.expectedKwargs), len(literal.KeywordArguments))
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

func TestMethodLiteralBlockArgument(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedArgs   []string
		expectedKwargs []string
		expectedBlock  string
	}{
		{
			name:          "with only a block argument",
			input:         "def method(&block); block; end",
			expectedArgs:  []string{},
			expectedBlock: "block",
		},
		{
			name:          "with a positional and a block argument",
			input:         "def method(value, &block); block; end",
			expectedArgs:  []string{"value"},
			expectedBlock: "block",
		},
		{
			name:           "with a keyword and a block argument",
			input:          "def method(value:, &block); block; end",
			expectedKwargs: []string{"value"},
			expectedBlock:  "block",
		},
		{
			name:          "with a default and a block argument",
			input:         "def method(value = 1, &block); block; end",
			expectedArgs:  []string{"value"},
			expectedBlock: "block",
		},
		{
			name:          "without parentheses after a positional argument",
			input:         "def method value, &block; block; end",
			expectedArgs:  []string{"value"},
			expectedBlock: "block",
		},
		{
			name:          "without parentheses and only a block argument",
			input:         "def method &block; block; end",
			expectedBlock: "block",
		},
		{
			name:          "after a line break",
			input:         "def method value,\n\t&block\n\tblock\nend",
			expectedArgs:  []string{"value"},
			expectedBlock: "block",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			testExpressionStatement(t, program.Statements[0], func(literal *ast.MethodLiteral) {
				if len(literal.Arguments) != len(tt.expectedArgs) {
					t.Fatalf("exp %d parameters got=%d", len(tt.expectedArgs), len(literal.Arguments))
				}

				for i, parameter := range literal.Arguments {
					testIdentifier(t, *parameter, tt.expectedArgs[i])
				}

				if len(literal.KeywordArguments) != len(tt.expectedKwargs) {
					t.Fatalf("exp %d keyword parameters got=%d", len(tt.expectedKwargs), len(literal.KeywordArguments))
				}

				for i, argument := range literal.KeywordArguments {
					testIdentifier(t, *argument, tt.expectedKwargs[i])
				}

				if literal.BlockArgument == nil {
					t.Fatalf("expected the block argument %q, got none", tt.expectedBlock)
				}

				testIdentifier(t, *literal.BlockArgument, tt.expectedBlock)
			})
		})
	}
}

func TestMethodLiteralParameterDefaults(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedArgs     []string
		expectedDefaults []string // empty string for a required parameter
		expectedKwargs   []string
	}{
		{
			name:             "trailing default",
			input:            "def method(a, b = 1); a; end",
			expectedArgs:     []string{"a", "b"},
			expectedDefaults: []string{"", "1"},
		},
		{
			name:             "all parameters defaulted",
			input:            "def method(a = 1, b = 2); a; end",
			expectedArgs:     []string{"a", "b"},
			expectedDefaults: []string{"1", "2"},
		},
		{
			name:             "default before a required parameter",
			input:            "def method(a = 1, b); b; end",
			expectedArgs:     []string{"a", "b"},
			expectedDefaults: []string{"1", ""},
		},
		{
			name:             "expression default",
			input:            "def method(a = 1 + 2); a; end",
			expectedArgs:     []string{"a"},
			expectedDefaults: []string{"(1 + 2)"},
		},
		{
			name:             "default referring to an earlier parameter",
			input:            "def method(a, b = a + 1); b; end",
			expectedArgs:     []string{"a", "b"},
			expectedDefaults: []string{"", "(a + 1)"},
		},
		{
			name:             "default and keyword arguments",
			input:            "def method(a = 1, b:); b; end",
			expectedArgs:     []string{"a"},
			expectedDefaults: []string{"1"},
			expectedKwargs:   []string{"b"},
		},
		{
			name:             "defaults without parentheses",
			input:            "def method a = 1, b = 2; a; end",
			expectedArgs:     []string{"a", "b"},
			expectedDefaults: []string{"1", "2"},
		},
		{
			name: "defaults across lines without parentheses",
			input: `def method a,
				b = 2,
				c = 3
				a
			end`,
			expectedArgs:     []string{"a", "b", "c"},
			expectedDefaults: []string{"", "2", "3"},
		},
		{
			name:             "method without defaults",
			input:            "def method(a); a; end",
			expectedArgs:     []string{"a"},
			expectedDefaults: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			testExpressionStatement(t, program.Statements[0], func(literal *ast.MethodLiteral) {
				if len(literal.Arguments) != len(tt.expectedArgs) {
					t.Fatalf("exp %d parameters got=%d", len(tt.expectedArgs), len(literal.Arguments))
				}

				if len(literal.ArgumentDefaults) != len(literal.Arguments) {
					t.Fatalf(
						"expected a default entry for each of the %d parameters, got=%d",
						len(literal.Arguments),
						len(literal.ArgumentDefaults),
					)
				}

				for i, parameter := range literal.Arguments {
					testIdentifier(t, *parameter, tt.expectedArgs[i])

					defaultValue := ""
					if defaultValueExpression := literal.ArgumentDefaults[i]; defaultValueExpression != nil {
						defaultValue = defaultValueExpression.String(0)
					}

					if defaultValue != tt.expectedDefaults[i] {
						t.Errorf("parameter %s default: want=%q got=%q", parameter.Value, tt.expectedDefaults[i], defaultValue)
					}
				}

				if len(literal.KeywordArguments) != len(tt.expectedKwargs) {
					t.Fatalf("exp %d keyword parameters got=%d", len(tt.expectedKwargs), len(literal.KeywordArguments))
				}

				for i, argument := range literal.KeywordArguments {
					testIdentifier(t, *argument, tt.expectedKwargs[i])
				}
			})
		})
	}
}

func TestMethodLiteralBodyWithGroupedExpression(t *testing.T) {
	program := testParseAST(t, `
		def foo
			(1 + 2) * 3
		end
	`)

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(literal *ast.MethodLiteral) {
		testLiteralExpression(t, literal.Name, "foo")
		expectStatementLength(t, literal.Body.Statements, 1)

		testExpressionStatement(t, literal.Body.Statements[0], func(body ast.Expression) {
			multiplied, ok := body.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("body is not *ast.InfixExpression. got=%T", body)
			}

			if multiplied.Operator != "*" {
				t.Fatalf("operator is not '*'. got=%q", multiplied.Operator)
			}

			testInfixExpression(t, multiplied.Left, 1, "+", 2)
			testLiteralExpression(t, multiplied.Right, 3)
		})
	})
}
