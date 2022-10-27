package compiler

import (
	ast "emerald/parser/ast"
	"fmt"
)

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) {
	var op Opcode

	switch node.Operator {
	case "+":
		op = OpAdd
	case "-":
		op = OpSub
	case "*":
		op = OpMul
	case "/":
		op = OpDiv
	case "=~":
		op = OpMatch
	case "<=>":
		op = OpSpaceship
	case ">":
		op = OpGreaterThan
	case ">=":
		op = OpGreaterThanOrEq
	case "<":
		op = OpLessThan
	case "<=":
		op = OpLessThanOrEq
	case "==":
		op = OpEqual
	case "===":
		op = OpCaseEqual
	case "!=":
		op = OpNotEqual
	case "&&":
		c.compileIfExpression(&ast.IfExpression{
			Condition:   node.Left,
			Consequence: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Right}}},
			Alternative: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Left}}},
		})
		return
	case "||":
		c.compileIfExpression(&ast.IfExpression{
			Condition:   node.Left,
			Consequence: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Left}}},
			Alternative: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Right}}},
		})
		return
	default:
		panic(fmt.Errorf("unknown infix operator %s", node.Operator))
	}

	c.Compile(node.Right)
	c.Compile(node.Left)
	c.emit(op)
}
