package parser

import (
	"emerald/ast"
	"testing"
)

func TestHashLiteralParsing(t *testing.T) {
	input := `{
		hello: false,
		0: true,
		"string": 25,
		nested: {
			key: 0,
		},
	}`

	program := testParseAST(t, input)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	value, ok := literal.Value["hello"]
	if !ok {
		t.Errorf("literal.Value doesn't have expected key 'hello'. got=%q", literal.Value)
	}

	testBooleanLiteral(t, value, false)

	value, ok = literal.Value["0"]
	if !ok {
		t.Errorf("literal.Value doesn't have expected key '0'. got=%q", literal.Value)
	}

	testBooleanLiteral(t, value, true)

	value, ok = literal.Value["string"]
	if !ok {
		t.Errorf("literal.Value doesn't have expected key 'string'. got=%q", literal.Value)
	}

	testIntegerLiteral(t, value, 25)

	value, ok = literal.Value["nested"]
	if !ok {
		t.Errorf("literal.Value doesn't have expected key 'nested'. got=%q", literal.Value)
	}

	hashValue, ok := value.(*ast.HashLiteral)
	if !ok {
		t.Errorf("literal.Value['nested'] was not hash got=%T", value)
	}

	testIntegerLiteral(t, hashValue.Value["key"], 0)
}
