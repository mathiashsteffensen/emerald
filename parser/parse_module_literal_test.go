package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestParseModuleLiteral(t *testing.T) {
	input := `
		module MyMod
			module NestedMod; end
			def do_stuff; end
		end
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(mod *ast.ModuleLiteral) {
		testIdentifier(t, mod.Name, "MyMod")

		expectStatementLength(t, mod.Body.Statements, 2)

		testExpressionStatement(t, mod.Body.Statements[0], func(mod *ast.ModuleLiteral) {
			testIdentifier(t, mod.Name, "NestedMod")
		})
		testExpressionStatement(t, mod.Body.Statements[1], func(method *ast.MethodLiteral) {
			testIdentifier(t, method.Name, "do_stuff")
		})
	})
}
