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

		if len(program.Statements) != 1 {
			t.Fatalf("program.Body does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T (%+v)", stmt.Expression, stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, tt.condition.left, tt.condition.op, tt.condition.right) {
			return
		}

		if len(exp.Consequence.Statements) != 1 {
			t.Errorf("consequence is not 1 statements. got=%d\n",
				len(exp.Consequence.Statements))
		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
		}

		if !testLiteralExpression(t, consequence.Expression, tt.consequence) {
			return
		}

		if exp.Alternative != nil {
			t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
		}
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

			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Body does not contain %d statements. got=%d\n (%+v)",
					1,
					len(program.Statements),
					program.Statements,
				)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.IfExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
			}

			if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
				return
			}

			if len(exp.Consequence.Statements) != 1 {
				t.Fatalf(
					"consequence is not 1 statements. got=%d (%+v)\n",
					len(exp.Consequence.Statements),
					exp.Consequence.Statements,
				)
			}

			consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
					exp.Consequence.Statements[0])
			}

			tt.expect(t, consequence.Expression)

			if len(exp.Alternative.Statements) != 1 {
				t.Fatalf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
					len(exp.Alternative.Statements))
			}

			alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
					exp.Alternative.Statements[0])
			}

			if !testIdentifier(t, alternative.Expression, "y") {
				return
			}
		})
	}
}
