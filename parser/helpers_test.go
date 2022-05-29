package parser

import (
	"emerald/ast"
	"emerald/lexer"
	"fmt"
	"testing"
)

func testParseAST(t *testing.T, input string) *ast.AST {
	l := lexer.New(lexer.NewInput("test.rb", input))

	p := New(l)
	program := p.ParseAST()

	checkParserErrors(t, p)

	return program
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.IdentifierExpression)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testRescueBlock(actual *ast.RescueBlock, expStmts int, expErrVarName string, expErrClasses ...string) error {
	if len(actual.CaughtErrorClasses) != len(expErrClasses) {
		return fmt.Errorf(
			"rescue block caught wrong amount of error classes \nwant=%d\ngot=%d",
			len(expErrClasses),
			len(actual.CaughtErrorClasses),
		)
	}

	for i, class := range actual.CaughtErrorClasses {
		if class.String() != expErrClasses[i] {
			return fmt.Errorf("CaughtErrorClasses[%d] failed \nwant=%s\ngot=%s", i, expErrClasses[i], class)
		}
	}

	if actual.ErrorVarName.String() != expErrVarName {
		return fmt.Errorf(
			"wrong rescue block error var name \nwant='%s'\ngot='%s'",
			expErrVarName,
			actual.ErrorVarName,
		)
	}

	if len(actual.Body.Statements) != expStmts {
		return fmt.Errorf(
			"rescue block had wrong amount of statements \nwant=%d\ngot=%d",
			expStmts,
			len(actual.Body.Statements),
		)
	}

	return nil
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser_test has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser_test error: %q", msg)
	}
	t.FailNow()
}
