package compiler

import (
	"emerald/ast"
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
		return c.compileBooleanAndOperator(node)
	case "||":
		return c.compileBooleanOrOperator(node)
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

func (c *Compiler) compileBooleanAndOperator(node *ast.InfixExpression) error {
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	// Emit an `OpJumpNotTruthy` with a bogus value
	jumpNotTruthyPos := c.emit(OpJumpNotTruthy, 9999)
	c.emit(OpPop)

	err = c.Compile(node.Right)
	if err != nil {
		return err
	}

	afterConsequencePos := len(c.currentInstructions()) + 1
	c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	return nil
}

func (c *Compiler) compileBooleanOrOperator(node *ast.InfixExpression) error {
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	// Emit an `OpJumpTruthy` with a bogus value
	jumpTruthyPos := c.emit(OpJumpTruthy, 9999)
	c.emit(OpPop)

	err = c.Compile(node.Right)
	if err != nil {
		return err
	}

	afterConsequencePos := len(c.currentInstructions()) + 1
	c.changeOperand(jumpTruthyPos, afterConsequencePos)

	return nil
}
