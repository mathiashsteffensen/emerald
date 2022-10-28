package parser

import (
	"emerald/parser/ast"
	"testing"
)

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{`5 + 5`, 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"false && true", false, "&&", true},
		{"false || true", false, "||", true},
		{"1 <=> 2", 1, "<=>", 2},
		{"regex =~ string", "regex", "=~", "string"},
		{"Integer === 1", "Integer", "===", 1},
	}

	for _, tt := range infixTests {
		program := testParseAST(t, tt.input)

		expectStatementLength(t, program.Statements, 1)

		testExpressionStatement(t, program.Statements[0], func(expression *ast.InfixExpression) {
			testInfixExpression(t, expression, tt.leftValue, tt.operator, tt.rightValue)
		})
	}
}
