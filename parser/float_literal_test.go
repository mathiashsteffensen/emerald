package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestFloatLiteral(t *testing.T) {
	input := "15.25; 4e10; 4.0E-1"

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 3)

	testExpressionStatement(t, program.Statements[0], func(float *ast.FloatLiteral) {
		testFloatLiteral(t, float, 15.25)
	})

	testExpressionStatement(t, program.Statements[1], func(float *ast.FloatLiteral) {
		testFloatLiteral(t, float, 4e10)
	})

	testExpressionStatement(t, program.Statements[2], func(float *ast.FloatLiteral) {
		testFloatLiteral(t, float, 4.0e-1)
	})
}
