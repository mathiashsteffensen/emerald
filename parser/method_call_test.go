package parser

import (
	"emerald/ast"
	"testing"
)

func TestMethodCallParsing(t *testing.T) {
	input := `
		1.add(2, 3, 4 + 5) do |num, next|
			num + @var.method
		rescue NoMemoryError, SystemError => e
			puts("we done fucked up this time")
		rescue Exception => e
			puts("Adding is hard :(")
		ensure
			puts("This will always run")
		end.first

		Logger.info(msg)
		Logger.level = :debug
	`

	program := testParseAST(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d (%+v)\n", len(program.Statements), program.Statements)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.MethodCall)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.MethodCall. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Method, "first") {
		return
	}

	exp, ok = exp.Left.(*ast.MethodCall)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.MethodCall. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Method, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 2)
	testLiteralExpression(t, exp.Arguments[1], 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)

	if exp.Block == nil {
		t.Fatalf("method call was not passed a block")
	}

	if len(exp.Block.RescueBlocks) != 2 {
		t.Fatalf("block was not passed 2 rescue clauses, got=%d", len(exp.Block.RescueBlocks))
	}

	rescue := exp.Block.RescueBlocks[0]

	err := testRescueBlock(rescue, 1, "e", "NoMemoryError", "SystemError")
	if err != nil {
		t.Errorf("first rescue failed: %s", err)
	}

	rescue = exp.Block.RescueBlocks[1]

	err = testRescueBlock(rescue, 1, "e", "Exception")
	if err != nil {
		t.Errorf("first rescue failed: %s", err)
	}

	if exp.Block.EnsureBlock == nil {
		t.Fatalf("block was not passed an ensure clause")
	}

	if len(exp.Block.EnsureBlock.Body.Statements) != 1 {
		t.Fatalf("rescue block did not have 1 statement got=%d", len(exp.Block.EnsureBlock.Body.Statements))
	}

	stmt, ok = program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok = stmt.Expression.(*ast.MethodCall)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.MethodCall. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Left, "Logger") {
		return
	}
	if !testIdentifier(t, exp.Method, "info") {
		return
	}

	stmt, ok = program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok = stmt.Expression.(*ast.MethodCall)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.MethodCall. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Left, "Logger") {
		return
	}
	if !testIdentifier(t, exp.Method, "level=") {
		return
	}
}
