package compiler

import (
	"emerald/ast"
	"fmt"
)

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) error {
	err := c.Compile(node.Right)
	if err != nil {
		return err
	}
	switch node.Operator {
	case "!":
		c.emit(OpBang)
	case "-":
		c.emit(OpMinus)
	default:
		return fmt.Errorf("unknown prefix operator %s", node.Operator)
	}

	return nil
}
