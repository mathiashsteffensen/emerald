package compiler

import (
	ast "emerald/parser/ast"
	"fmt"
)

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) error {
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
	case "!=":
		op = OpNotEqual
	case "&&":
		return c.compileIfExpression(&ast.IfExpression{
			Condition:   node.Left,
			Consequence: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Right}}},
			Alternative: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Left}}},
		})
	case "||":
		return c.compileIfExpression(&ast.IfExpression{
			Condition:   node.Left,
			Consequence: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Left}}},
			Alternative: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Right}}},
		})
	case "&&=":
		assignment := &ast.AssignmentExpression{
			Name:  node.Left,
			Value: node.Right,
		}
		return c.compileIfExpression(&ast.IfExpression{
			Condition:   node.Left,
			Consequence: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: assignment}}},
			Alternative: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Left}}},
		})
	case "||=":
		assignment := &ast.AssignmentExpression{
			Name:  node.Left,
			Value: node.Right,
		}
		return c.compileIfExpression(&ast.IfExpression{
			Condition:   node.Left,
			Consequence: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: node.Left}}},
			Alternative: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: assignment}}},
		})
	default:
		return fmt.Errorf("unknown infix operator %s", node.Operator)
	}

	err := c.Compile(node.Right)
	if err != nil {
		return err
	}

	err = c.Compile(node.Left)
	if err != nil {
		return err
	}

	c.emit(op)

	return nil
}
