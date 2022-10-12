package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestClassLiteral(t *testing.T) {
	input := `class Math < BasicObject
		def add(x, y)
			x + y
		end
	end`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(class *ast.ClassLiteral) {
		if class.Name.String(0) != "Math" {
			t.Fatalf("class name was not Math got=%s", class.Name.String(0))
		}

		if class.Parent.String(0) != "BasicObject" {
			t.Fatalf("parent class name was not BasicObject got=%s", class.Name.String(0))
		}

		expectStatementLength(t, class.Body.Statements, 1)

		testExpressionStatement(t, class.Body.Statements[0], func(method *ast.MethodLiteral) {
			if method.Name.String(0) != "add" {
				t.Fatalf("method name was not add got=%s", method.Name.String(0))
			}

			expectStatementLength(t, method.Body.Statements, 1)
		})
	})
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

	if class.Name.String(0) != "Logger" {
		t.Fatalf("class name was not integer got=%s", class.Name.String(0))
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

	if method.Name.String(0) != "info" {
		t.Fatalf("method name was not add got=%s", method.Name.String(0))
	}

	if len(method.Body.Statements) != 1 {
		t.Fatalf(
			"expected method body to have 1 statement got=%d (%+v)",
			len(class.Body.Statements),
			class.Body.Statements,
		)
	}
}
