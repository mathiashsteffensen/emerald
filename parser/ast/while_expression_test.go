package ast_test

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"testing"
)

func TestWhileExpression(t *testing.T) {
	tests := []struct {
		token    lexer.Token
		isUntil  bool
		expected string
	}{
		{token: lexer.Token{Type: lexer.WHILE, Literal: "while"}, expected: "while running\n\trunning\nend"},
		{token: lexer.Token{Type: lexer.UNTIL, Literal: "until"}, isUntil: true, expected: "until running\n\trunning\nend"},
	}

	for _, tt := range tests {
		expression := &ast.WhileExpression{
			Token:     tt.token,
			Condition: ast.IdentifierExpression{Value: "running"},
			Consequence: &ast.BlockStatement{Statements: []ast.Statement{
				&ast.ExpressionStatement{Expression: ast.IdentifierExpression{Value: "running"}},
			}},
		}

		if expression.TokenLiteral() != tt.token.Literal {
			t.Errorf("wrong token literal, got=%q", expression.TokenLiteral())
		}

		if expression.IsUntil() != tt.isUntil {
			t.Errorf("wrong IsUntil, want=%v got=%v", tt.isUntil, expression.IsUntil())
		}

		if expression.String(0) != tt.expected {
			t.Errorf("wrong string, want=%q got=%q", tt.expected, expression.String(0))
		}
	}
}
