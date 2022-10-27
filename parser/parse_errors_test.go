package parser

import "testing"

func TestParseErrors(t *testing.T) {
	testParseError(t, "def method", "syntax error, unexpected end-of-input")
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
}
