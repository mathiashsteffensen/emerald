package parser

import (
	"emerald/ast"
	"testing"
)

func TestFloatLiteral(t *testing.T) {
	input := "15.25"

	program := testParseAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	expr, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected expression got=%#v", program.Statements[0])
	}

	float, ok := expr.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected float got=%#v", expr.Expression)
	}

	if float.Value != 15.25 {
		t.Fatalf("expected float to have value 15.25, got=%f", float.Value)
	}

	if float.String() != "15.25" {
		t.Fatalf("expected float to have value 15.25, got=%s", float.String())
	}
}
