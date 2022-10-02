package parser

import "testing"

func TestParseErrors(t *testing.T) {
	testParseError(t, "def method", "syntax error, unexpected end-of-input")
	testParseError(t, ":2", "expected next token to be one of [IDENT, STRING], got INT instead")
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
	testParseError(t, "(2+2", "expected next token to be ), got EOF instead")
	testParseError(t, "p(", "expected next token to be ), got EOF instead")
	testParseError(t, "p { 2", "expected next token to be }, got EOF instead")
	testParseError(t, `"hello #{name"`, "expected next token to be }, got EOF instead")
}
