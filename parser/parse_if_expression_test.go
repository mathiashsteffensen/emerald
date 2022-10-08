package parser

import (
	ast "emerald/parser/ast"
	"testing"
)

func TestIfExpression(t *testing.T) {
	type conditionTest struct {
		left  any
		op    string
		right any
	}

	type test struct {
		input       string
		condition   conditionTest
		consequence any
	}

	tests := []test{
		{
			`
			if x < y
				x
			end
			`,
			conditionTest{"x", "<", "y"},
			"x",
		},
		{
			"if 5 < 98; false; end",
			conditionTest{5, "<", 98},
			false,
		},
		{
			"false if 5 < 98",
			conditionTest{5, "<", 98},
			false,
		},
		{
			"true if 10 * 5",
			conditionTest{10, "*", 5},
			true,
		},
	}

	for _, tt := range tests {
		program := testParseAST(t, tt.input)

		expectStatementLength(t, program.Statements, 1)

		testExpressionStatement(t, program.Statements[0], func(exp *ast.IfExpression) {
			testInfixExpression(t, exp.Condition, tt.condition.left, tt.condition.op, tt.condition.right)

			expectStatementLength(t, exp.Consequence.Statements, 1)

			testExpressionStatement(t, exp.Consequence.Statements[0], func(expression ast.Expression) {
				testLiteralExpression(t, expression, tt.consequence)
				if exp.Alternative != nil {
					t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
				}
			})
		})
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect func(t *testing.T, evaluated ast.Expression)
	}{
		{
			"one-liner",
			"if x < y;  x; else; y; end;",
			func(t *testing.T, evaluated ast.Expression) {
				if !testIdentifier(t, evaluated, "x") {
					return
				}
			},
		},
		{
			"multi line",
			`if x < y
				x
			 else
				y
			end`,
			func(t *testing.T, evaluated ast.Expression) {
				if !testIdentifier(t, evaluated, "x") {
					return
				}
			},
		},
		{
			"multi line with method call",
			`if x < y
				Logger.info("Hello World!")
			 else
				y
			end`,
			func(t *testing.T, evaluated ast.Expression) {
				if _, ok := evaluated.(*ast.MethodCall); !ok {
					t.Fatalf("exp not method call got=%T", evaluated)
				}
			},
		},
		{
			"multi line and nested",
			`if x < y
				if x
					x
				end
			 else
				y
			end`,
			func(t *testing.T, evaluated ast.Expression) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			testExpressionStatement(t, program.Statements[0], func(expression *ast.IfExpression) {
				testInfixExpression(t, expression.Condition, "x", "<", "y")

				expectStatementLength(t, expression.Consequence.Statements, 1)
				testExpressionStatement(t, expression.Consequence.Statements[0], func(expression ast.Expression) {
					tt.expect(t, expression)
				})

				expectStatementLength(t, expression.Alternative.Statements, 1)
				testExpressionStatement(t, expression.Alternative.Statements[0], func(expression ast.Expression) {
					testIdentifier(t, expression, "y")
				})
			})
		})
	}
}

func TestIfElsifExpression(t *testing.T) {
	input := `
		if x < y
			a
		elsif true
			b
		elsif 2 + 2 == 4 
			c
		end
	`

	program := testParseAST(t, input)

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(expression *ast.IfExpression) {
		testInfixExpression(t, expression.Condition, "x", "<", "y")

		if len(expression.ElseIfs) != 2 {
			t.Fatalf("Expected 2 elsif statements but got %d", len(expression.ElseIfs))
		}

		testLiteralExpression(t, expression.ElseIfs[0].Condition, true)
		if expression.ElseIfs[1].Condition.String() != "((2 + 2) == 4)" {
			t.Errorf("Expected elsif condition to be ((2 + 2) == 4) but got %s", expression.ElseIfs[1].Condition.String())
		}
	})
}

func TestUnlessExpression(t *testing.T) {
	type conditionTest struct {
		left  any
		op    string
		right any
	}

	type test struct {
		input       string
		condition   conditionTest
		alternative any
	}

	tests := []test{
		{
			`
			unless x < y
				x
			end
			`,
			conditionTest{"x", "<", "y"},
			"x",
		},
		{
			"unless 5 < 98; false; end",
			conditionTest{5, "<", 98},
			false,
		},
		{
			"false unless 5 < 98",
			conditionTest{5, "<", 98},
			false,
		},
		{
			"true unless 10 * 5",
			conditionTest{10, "*", 5},
			true,
		},
	}

	for _, tt := range tests {
		program := testParseAST(t, tt.input)

		expectStatementLength(t, program.Statements, 1)

		testExpressionStatement(t, program.Statements[0], func(exp *ast.IfExpression) {
			testInfixExpression(t, exp.Condition, tt.condition.left, tt.condition.op, tt.condition.right)
			if exp.Consequence != nil {
				t.Errorf("exp.Consequence was not nil. got=%+v", exp.Consequence)
			}

			expectStatementLength(t, exp.Alternative.Statements, 1)

			testExpressionStatement(t, exp.Alternative.Statements[0], func(expression ast.Expression) {
				testLiteralExpression(t, expression, tt.alternative)
			})
		})
	}
}

func TestUnlessElseExpression(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect func(t *testing.T, evaluated ast.Expression)
	}{
		{
			"one-liner",
			"unless x < y;  x; else; y; end;",
			func(t *testing.T, evaluated ast.Expression) {
				testIdentifier(t, evaluated, "x")
			},
		},
		{
			"multi line",
			`unless x < y
				x
			 else
				y
			end`,
			func(t *testing.T, evaluated ast.Expression) {
				testIdentifier(t, evaluated, "x")
			},
		},
		{
			"multi line with method call",
			`unless x < y
				Logger.info("Hello World!")
			 else
				y
			end`,
			func(t *testing.T, evaluated ast.Expression) {
				if _, ok := evaluated.(*ast.MethodCall); !ok {
					t.Fatalf("exp not method call got=%T", evaluated)
				}
			},
		},
		{
			"multi line and nested",
			`unless x < y
				unless x
					x
				end
			 else
				y
			end`,
			func(t *testing.T, evaluated ast.Expression) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := testParseAST(t, tt.input)

			expectStatementLength(t, program.Statements, 1)

			testExpressionStatement(t, program.Statements[0], func(expression *ast.IfExpression) {
				testInfixExpression(t, expression.Condition, "x", "<", "y")

				expectStatementLength(t, expression.Alternative.Statements, 1)
				testExpressionStatement(t, expression.Alternative.Statements[0], func(expression ast.Expression) {
					tt.expect(t, expression)
				})

				expectStatementLength(t, expression.Consequence.Statements, 1)
				testExpressionStatement(t, expression.Consequence.Statements[0], func(expression ast.Expression) {
					testIdentifier(t, expression, "y")
				})
			})
		})
	}
}
