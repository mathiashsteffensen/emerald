package parser

import (
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"strings"
	"testing"
)

func testParseAST(t *testing.T, input string) *ast.AST {
	l := lexer.New(lexer.NewInput("test.rb", input))

	p := New(l)
	program := p.ParseAST()

	checkParserErrors(t, p)

	return program
}

func testParseError(t *testing.T, input, expectedError string) {
	l := lexer.New(lexer.NewInput("test.rb", input))

	p := New(l)
	p.ParseAST()

	if len(p.Errors()) == 0 {
		t.Fatalf("Expected parser to have errors, but had none")
	}

	if p.Errors()[0] != expectedError {
		t.Fatalf("Expected parser to have error=%q but got=%q", expectedError, p.Errors()[0])
	}
}

func expectStatementLength(t *testing.T, stmt []ast.Statement, length int) {
	if len(stmt) != length {
		t.Fatalf("AST does not contain %d statements. got=%d\n\n%s", length, len(stmt), (&ast.BlockStatement{Statements: stmt}).String())
	}
}

func testIfElseExpression(t *testing.T, expr *ast.IfExpression, condition, consequence, alternative string) {
	if strings.Trim(expr.Condition.String(), "\n\r\t ") != condition {
		t.Errorf("Unexpected condition want=%s got=%s", condition, expr.Condition.String())
	}

	if strings.Trim(expr.Consequence.String(), "\n\r\t ") != consequence {
		t.Errorf("Unexpected consequence want=%s got=%s", consequence, expr.Consequence.String())
	}

	if strings.Trim(expr.Alternative.String(), "\n\r\t ") != alternative {
		t.Errorf("Unexpected alternative want=%s got=%s", alternative, expr.Alternative.String())
	}
}

func testBreakStatement(t *testing.T, stmt ast.Statement, value any) {
	breakStmt, ok := stmt.(*ast.BreakStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.BreakStatement. got=%T", stmt)
	}

	if value == nil && breakStmt.Value != nil {
		t.Fatalf("Expected break without a value but got %s", breakStmt.Value)
	} else if value == nil {
		return
	}

	testLiteralExpression(t, breakStmt.Value, value)
}

func testExpressionStatement[T ast.Expression](t *testing.T, stmt ast.Statement, cb func(expression T)) {
	exprStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
	}

	testExpression(t, exprStmt.Expression, cb)
}

func testExpression[T ast.Expression](t *testing.T, expr ast.Expression, cb func(expression T)) {
	exp, ok := expr.(T)
	if !ok {
		t.Fatalf("Expression is not expected type, got=%T", expr)
	}

	cb(exp)
}

func testAssignmentExpression(t *testing.T, expr ast.Expression, expectedName, expectedValue string) {
	ident, ok := expr.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("exp not *ast.AssignmentExpression. got=%T", expr)
	}
	if ident.Name.String() != expectedName {
		t.Errorf("ident.Name not %s. got=%s", expectedName, ident.Name.String())
	}
	if ident.Value.String() != expectedValue {
		t.Errorf("ident.Value not %s. got=%s", expectedValue, ident.Value.String())
	}
}

func testCallExpression(t *testing.T, expr ast.Expression, name string, args []any, block bool) {
	exp, ok := expr.(ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", expr)
	}

	if exp.Method.String() != name {
		t.Fatalf("Method name was not %s, got=%s", name, exp.Method.String())
	}

	if len(exp.Arguments) != len(args) {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	for i, arg := range args {
		testLiteralExpression(t, exp.Arguments[i], arg)
	}

	if block && exp.Block == nil {
		t.Fatalf("exp was not passed a block")
	}
}

func testMethodCall(t *testing.T, expr ast.Expression, receiver string, name string, args []any, block bool) {
	exp, ok := expr.(*ast.MethodCall)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.MethodCall. got=%T", expr)
	}

	testIdentifier(t, exp.Left, receiver)
	testCallExpression(t, exp.CallExpression, name, args, block)
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
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
	expected any,
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		if strings.HasPrefix(v, ":") {
			return testSymbolLiteral(t, exp, v[1:])
		}
		if strings.HasPrefix(v, "s:") {
			return testStringLiteral(t, exp, v[2:])
		}
		if strings.HasPrefix(v, "@") {
			return testInstanceVar(t, exp, v)
		}
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testSymbolLiteral(t *testing.T, expr ast.Expression, expected string) bool {
	sym, ok := expr.(*ast.SymbolLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.SymbolLiteral. got=%T", expr)
	}

	if sym.Value != expected {
		t.Errorf("sym.Value not %s. got=%s", expected, sym.Value)
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, expr ast.Expression, expected string) bool {
	str, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.StringLiteral. got=%T", expr)
	}

	if str.Value != expected {
		t.Errorf("str.Value not %q. got=%q", expected, str.Value)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value,
			integer.TokenLiteral())
		return false
	}

	return true
}

func testInstanceVar(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.InstanceVariable)
	if !ok {
		t.Errorf("exp not *ast.InstanceVariable. got=%T", exp)
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(ast.IdentifierExpression)
	if !ok {
		t.Errorf("exp not *ast.Identifier(%s). got=%T", value, exp)
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

func testHashLiteral(t *testing.T, expr ast.Expression, items map[string]any) {
	hash, ok := expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("Expected HashLiteral but got %T", expr)
	}

	for expectedKey, expectedValue := range items {
		keyFound := false

		for actualKey, actualValue := range hash.Value {
			if actualKey.String() == expectedKey {
				keyFound = true
				testLiteralExpression(t, actualValue, expectedValue)
			}
		}

		if !keyFound {
			t.Fatalf("HashLiteral did not have expected key %q, hash %#v", expectedKey, hash.Value)
		}
	}
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
