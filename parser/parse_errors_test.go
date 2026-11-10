package parser

import "testing"

func TestParseErrors(t *testing.T) {
	testParseError(t, "def method", "syntax error, unexpected end-of-input")
	testParseError(t, "def method(a", "syntax error, unexpected end-of-input")
	testParseError(t, "{ hello }", "expected next token to be =>, got } instead")
	testParseError(
		t,
		`
				{
					hello: 2
					HELLO: 3
				}
		`,
		"expected next token to be }, got IDENT instead",
	)
	testParseError(t, "(2+2", "syntax error, unexpected end-of-input")
	testParseError(t, "p(", "syntax error, unexpected end-of-input")
	testParseError(t, "p { 2", "syntax error, unexpected end-of-input")
	testParseError(t, `"hello #{name"`, "syntax error, unexpected end-of-input")
	testParseError(t, `true ? 2`, "syntax error, unexpected end-of-input")
	testParseError(t, "def method(, a); end", "expected next token to be IDENT")
	testParseError(t, "def method(a,.); end", "expected next token to be IDENT")
	testParseError(t, "def method(a:, b); end", "expected next token to be :")
	testParseError(t, "def method(&block, other); end", "expected next token to be ), got , instead")
	testParseError(t, "def method(&); end", "expected next token to be IDENT, got ) instead")
	testParseError(t, "def method a, &block, other; end", "expected a newline, ';' or end of input after the method parameters")
	testParseError(t, "def method &; end", "expected next token to be IDENT, got ; instead")
	testParseError(t, "def method &block, other; end", "expected a newline, ';' or end of input after the method parameters")
	testParseError(t, "def method a, &; end", "expected next token to be IDENT, got ; instead")
	testParseError(t, "def method(a, &a); end", "duplicated argument name")
	testParseError(t, "def method(a = 1, &a); end", "duplicated argument name")
	testParseError(t, "def method(a:, &a); end", "duplicated argument name")
	testParseError(t, "def method a, &a; end", "duplicated argument name")
	testParseError(t, "def method a:, &a; end", "duplicated argument name")
	testParseError(t, "def method a b; end", "expected a newline, ';' or end of input after the method parameters at line 1, column 13\ndef method a b; end\n             ^")
	testParseError(t, "def method a,; end", "expected next token to be IDENT")
	testParseError(t, "def method a:, b; end", "expected next token to be :")
	testParseError(t, "def foo\n  (a, b)\nend", "expected next token to be ), got , instead")
	testParseError(t, `{"x":}`, `expected a value for hash key "x"`)
	testParseError(t, "{valid?:}", `expected a value for hash key "valid?"`)
	testParseError(t, "{@x:}", `expected a value for hash key "@x"`)
	testParseError(t, "h = {\n  x:", "syntax error, unexpected end-of-input")
	testParseError(t, "add(x:", "syntax error, unexpected end-of-input")
	testParseError(t, "def method(for); end", "expected next token to be :, got ) instead")
	testParseError(t, "for in [1]; end", "expected next token to be IDENT, got IN instead at line 1 column 4")
	testParseError(t, "for x [1]; end", "expected next token to be IN, got [ instead at line 1 column 6")
	testParseError(t, "for a, b in [[1, 2]]; end", "expected next token to be IN, got , instead")
	testParseError(t, "for x in [1]; puts(x)", "syntax error, unexpected end-of-input")
	testParseError(t, "for x in [1]", "syntax error, unexpected end-of-input")
	testParseError(t, "for x in", "syntax error, unexpected end-of-input")
	testParseError(t, "for x in [1] puts(x) end", "expected 'do', a newline or ';' after the for loop collection at line 1, column 13")
}

func TestParseErrorsOnUnparsableExpressionDoNotHangOrPanic(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{"?.foo", `no prefix parse function for ?:"?"`},
		{"?::Foo", `no prefix parse function for ?:"?"`},
		{"? >= 1", `no prefix parse function for ?:"?"`},
		{"? && 1", `no prefix parse function for ?:"?"`},
		{"? || 1", `no prefix parse function for ?:"?"`},
		{")\n.foo", `no prefix parse function for ):")"`},
		{"\xC3.foo", `invalid multibyte char (UTF-8) \xC3 at line 1, column 0`},
		{"\xC3(1)", `invalid multibyte char (UTF-8) \xC3 at line 1, column 0`},
		{"1(2)", "syntax error, unexpected '(' at line 1, column 1"},
		{`"a"(1)`, "syntax error, unexpected '('"},
		{"@a(1)", "syntax error, unexpected '('"},
		{"\x01(", `no prefix parse function for ILLEGAL:"\x01" found at line 1, column 0`},
		{"class X = 1", "class/module name must be CONSTANT at line 1, column 6"},
		{"module x y", "class/module name must be CONSTANT at line 1, column 7"},
		// Only the line: `::` reports its second character's column, a known
		// fault of combineCurrentAndPeek.
		{"class :::Foo", "class/module name must be CONSTANT at line 1"},
		{"module \xF0\x9F", `invalid multibyte char (UTF-8) \xF0 at line 1, column 7`},
		{"class \xC3", `invalid multibyte char (UTF-8) \xC3 at line 1, column 6`},
		{"a = 1\nb = ?", `no prefix parse function for ?:"?" found at line 2, column 4`},
		{"a = 1\nb = 2\nπ = ?", `no prefix parse function for ?:"?" found at line 3, column 4`},
		{`"#{"` + "\xC3" + `"}"`, `invalid multibyte char (UTF-8) \xC3 at line 1, column 4`},
		{`x && "#{"` + "\xC3" + `"}"`, `invalid multibyte char (UTF-8) \xC3`},
		{"module || 1", "class/module name must be CONSTANT"},
		{"class && x", "class/module name must be CONSTANT"},
		{">= 1", `no prefix parse function for >=`},
	}

	for _, tt := range tests {
		testParseError(t, tt.input, tt.expectedError)
	}
}
