package compiler

import (
	"emerald/ast"
	"fmt"
)

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) error {
	err := c.Compile(node.Right)
	if err != nil {
		return err
	}

	err = c.Compile(node.Left)
	if err != nil {
		return err
	}

	switch node.Operator {
	case "+":
		c.emit(OpAdd)
	case "-":
		c.emit(OpSub)
	case "*":
		c.emit(OpMul)
	case "/":
		c.emit(OpDiv)
	case ">":
		c.emit(OpGreaterThan)
	case ">=":
		c.emit(OpGreaterThanOrEq)
	case "<":
		c.emit(OpLessThan)
	case "<=":
		c.emit(OpLessThanOrEq)
	case "==":
		c.emit(OpEqual)
	case "!=":
		c.emit(OpNotEqual)
	default:
		return fmt.Errorf("unknown infix operator %s", node.Operator)
	}

	return nil
}
