package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestParseModuleLiteral(t *testing.T) {
	input := `
		module MyMod
			def do_stuff; end
		end
	`

	program := testParseAST(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected %d statements, got %d", 1, len(program.Statements))
	}

	expr, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement got=%T", program.Statements[0])
	}

	mod, ok := expr.Expression.(*ast.ModuleLiteral)
	if !ok {
		t.Fatalf("expected *ast.ModuleLiteral got=%T", expr.Expression)
	}

	if mod.Name.String() != "MyMod" {
		t.Fatalf("expected name to be MyMod, but got=%s", mod.Name.String())
	}

	method := mod.Body.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.MethodLiteral)

	if method.Name.String() != "do_stuff" {
		t.Fatalf("expected name to be do_stuff, but got=%s", method.Name.String())
	}
}
