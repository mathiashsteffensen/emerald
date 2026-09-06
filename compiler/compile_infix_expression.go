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
	case "&&", "||":
		c.Compile(node.Left)
		c.emit(bytecode.OpDupN, node.Token, 1)
		jump := bytecode.OpJumpNotTruthy
		if node.Operator == "||" {
			jump = bytecode.OpJumpTruthy
		}
		end := c.emit(jump, node.Token, 9999)
		c.emit(bytecode.OpPop, node.Token)
		c.emit(bytecode.OpPop, node.Token)
		c.Compile(node.Right)
		c.changeOperand(end, len(c.currentInstructions()))
		return
	default:
		panic(fmt.Errorf("unknown infix operator %s", node.Operator))
	}

	c.Compile(node.Left)
	c.Compile(node.Right)
	c.emit(op, node.Token)
}
