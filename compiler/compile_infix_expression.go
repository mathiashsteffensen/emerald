package compiler

import (
	"emerald/bytecode"
	ast "emerald/parser/ast"
	"fmt"
)

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) {
	var op bytecode.Opcode

	switch node.Operator {
	case "+":
		op = bytecode.OpAdd
	case "-":
		op = bytecode.OpSub
	case "*":
		op = bytecode.OpMul
	case "/":
		op = bytecode.OpDiv
	case "=~":
		op = bytecode.OpMatch
	case "<=>":
		op = bytecode.OpSpaceship
	case ">":
		op = bytecode.OpGreaterThan
	case ">=":
		op = bytecode.OpGreaterThanOrEq
	case "<":
		op = bytecode.OpLessThan
	case "<=":
		op = bytecode.OpLessThanOrEq
	case "==":
		op = bytecode.OpEqual
	case "===":
		op = bytecode.OpCaseEqual
	case "!=":
		op = bytecode.OpNotEqual
	case "<<":
		op = bytecode.OpBinShiftLeft
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

	c.Compile(node.Left)
	c.Compile(node.Right)
	c.emit(op, node.Token)
}
