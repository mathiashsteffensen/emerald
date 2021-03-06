package parser

import (
	"emerald/ast"
	"testing"
)

func TestClassLiteral(t *testing.T) {
	input := `class Integer
		def add(x, y)
			x + y
		end
	end`

	program := testParseAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected program to have 1 statement got=%d", len(program.Statements))
	}

	class, ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.ClassLiteral)
	if !ok {
		t.Fatalf("expression is not class literal, got=%T", program.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if class.Name.String() != "Integer" {
		t.Fatalf("class name was not integer got=%s", class.Name.String())
	}

	if len(class.Body.Statements) != 1 {
		t.Fatalf(
			"expected class body to have 1 statement got=%d (%+v)",
			len(class.Body.Statements),
			class.Body.Statements,
		)
	}

	method, ok := class.Body.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.MethodLiteral)
	if !ok {
		t.Fatalf("expression is not method literal, got=%T", class.Body.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if method.Name.String() != "add" {
		t.Fatalf("method name was not add got=%s", method.Name.String())
	}

	if len(method.Body.Statements) != 1 {
		t.Fatalf(
			"expected method body to have 1 statement got=%d (%+v)",
			len(class.Body.Statements),
			class.Body.Statements,
		)
	}
}

func TestStaticClassLiteral(t *testing.T) {
	input := `class Logger
		class << self
			def info(msg)
				puts("INFO | " + msg)
			end
		end
	end`

	program := testParseAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected program to have 1 statement got=%d", len(program.Statements))
	}

	class, ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.ClassLiteral)
	if !ok {
		t.Fatalf("expression is not class literal, got=%T", program.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if class.Name.String() != "Logger" {
		t.Fatalf("class name was not integer got=%s", class.Name.String())
	}

	if len(class.Body.Statements) != 1 {
		t.Fatalf(
			"expected class body to have 1 statement got=%d (%+v)",
			len(class.Body.Statements),
			class.Body.Statements,
		)
	}

	static, ok := class.Body.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.StaticClassLiteral)
	if !ok {
		t.Fatalf("expression is not StaticClassLiteral, got=%T", class.Body.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if len(static.Body.Statements) != 1 {
		t.Fatalf(
			"expected class body to have 1 statement got=%d (%+v)",
			len(class.Body.Statements),
			class.Body.Statements,
		)
	}

	method, ok := static.Body.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.MethodLiteral)
	if !ok {
		t.Fatalf("expression is not method literal, got=%T", class.Body.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if method.Name.String() != "info" {
		t.Fatalf("method name was not add got=%s", method.Name.String())
	}

	if len(method.Body.Statements) != 1 {
		t.Fatalf(
			"expected method body to have 1 statement got=%d (%+v)",
			len(class.Body.Statements),
			class.Body.Statements,
		)
	}
}
