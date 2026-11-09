package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"strings"
	"testing"
)

func TestParseUnicodeIdentifiers(t *testing.T) {
	program := testParseAST(t, "café = 2\nnaïve?(ключ: 1)\nkw(ключ:)")

	expectStatementLength(t, program.Statements, 3)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.AssignmentExpression) {
		testAssignmentExpression(t, expression, "café", "2")
	})

	testExpressionStatement(t, program.Statements[1], func(expression ast.CallExpression) {
		testCallExpression(t, expression, "naïve?", []any{}, []string{":ключ"}, false)
	})

	testExpressionStatement(t, program.Statements[2], func(expression ast.CallExpression) {
		testCallExpression(t, expression, "kw", []any{}, []string{":ключ"}, false)
	})
}

// The lexer has no token for these ASCII operator characters and emits them as
// ILLEGAL; only invalid UTF-8 is an error, so they still parse as method names.
func TestParseOperatorCharactersWithoutTokenAsMethodNames(t *testing.T) {
	inputs := []string{
		"def %(o)\no\nend",
		"def ~\n7\nend",
		"def ^(o)\no\nend",
		"def `(c)\nc\nend",
		"x.%(4)",
		"x.^(1)",
		"x.~",
	}

	for _, input := range inputs {
		testParseAST(t, input)
	}
}

func TestParseInvalidUTF8IsNeverAName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		column int
	}{
		{"after def", "def \xC3\n42\nend", 4},
		{"after class", "class \xC3\nend", 6},
		{"after module", "module \xC3\nend", 7},
		{"after .", "x.\xC3", 2},
		{"after ::", "Foo::\xC3", 5},
		{"as hash label", "{\xC3: 1}", 1},
		{"as keyword argument", "m(\xC3: 1)", 2},
		{"after $", "$\xC3", 1},
		{"after @", "@\xC3", 1},
		{"after :", ":\xC3", 1},
		{"inside an array", "x = [1 \xC3]", 7},
		{"inside parentheses", "(1 \xC3", 3},
		{"as parameter", "def f(\xC3)\nend", 6},
		{"inside interpolation", "\"#{1 \xC3}\"", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedError := fmt.Sprintf(`invalid multibyte char (UTF-8) \xC3 at line 1, column %d`, tt.column)

			testParseError(t, tt.input, expectedError)
		})
	}
}

func TestParseInvalidUTF8InLiterals(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{"double-quoted string", "x = \"a\xC3b\"", `invalid multibyte char (UTF-8) \xC3 at line 1, column 6`},
		{"single-quoted string", "x = 'a\xC3b'", `invalid multibyte char (UTF-8) \xC3 at line 1, column 6`},
		{"string template", "y = 1\nx = \"a#{y}\xC3\"", `invalid multibyte char (UTF-8) \xC3 at line 2, column 10`},
		{"regexp", "x = /a\xC3b/", `invalid multibyte char (UTF-8) \xC3 at line 1, column 6`},
		{"quoted symbol", "x = :\"a\xC3\"", `invalid multibyte char (UTF-8) \xC3 at line 1, column 7`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParseError(t, tt.input, tt.expectedError)
		})
	}
}

func TestParseValidMultibyteStringAndInvalidUTF8InComment(t *testing.T) {
	program := testParseAST(t, "# a\xC3b\nx = \"é😀\"")

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.AssignmentExpression) {
		testAssignmentExpression(t, expression, "x", `"é😀"`)
	})
}

func TestParseInvalidByteIsReportedAfterTheErrorsInFrontOfIt(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{"foo(1,,\xC3)", `no prefix parse function for ,:"," found at line 1, column 6`},
		{")\xC3", `no prefix parse function for ):")" found at line 1, column 0`},
	}

	for _, tt := range tests {
		testParseError(t, tt.input, tt.expectedError)
	}
}

func TestParseReportsEachInvalidByteOnce(t *testing.T) {
	inputs := []string{
		"def \xC3\n42\nend",
		"class \xC3\nend",
		"module \xC3\nend",
		"x.\xC3",
		"Foo::\xC3",
		"{\xC3: 1}",
		"m(\xC3: 1)",
		"$\xC3",
		"@\xC3",
		":\xC3",
		"\xC3.foo",
		"foo(1,,\xC3)",
		")\xC3",
		"x = \"a\xC3\xC3\"",
		"module \xF0\x9F",
	}

	for _, input := range inputs {
		p := New(lexer.New(lexer.NewInput("test.rb", input)))
		p.ParseAST()

		reported := map[string]bool{}

		for _, err := range p.Errors() {
			if !strings.HasPrefix(err, "invalid multibyte char") {
				continue
			}

			location := strings.SplitN(err, "\n", 2)[0]
			if reported[location] {
				t.Errorf("%q: reported %q twice", input, location)
			}

			reported[location] = true
		}

		if len(reported) == 0 {
			t.Errorf("%q: expected an invalid multibyte char error, got %q", input, p.Errors())
		}
	}
}
