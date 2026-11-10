package ast_test

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"testing"
)

func TestForExpressionString(t *testing.T) {
	expression := &ast.ForExpression{
		Token:    lexer.Token{Literal: "for"},
		Variable: ast.IdentifierExpression{Value: "item"},
		Iterable: ast.IdentifierExpression{Value: "items"},
		Body: &ast.BlockStatement{Statements: []ast.Statement{
			&ast.ExpressionStatement{Expression: ast.IdentifierExpression{Value: "item"}},
		}},
	}

	if expression.TokenLiteral() != "for" {
		t.Errorf("wrong token literal, got=%q", expression.TokenLiteral())
	}

	expected := "for item in items\n\titem\nend"
	if expression.String(0) != expected {
		t.Errorf("wrong string, want=%q got=%q", expected, expression.String(0))
	}
}
