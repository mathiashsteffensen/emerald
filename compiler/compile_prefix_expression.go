package compiler

import (
	"emerald/parser/ast"
	"fmt"
)

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) {
	c.Compile(node.Right)

	switch node.Operator {
	case "!":
		c.emit(OpBang)
	case "-":
		c.emit(OpMinus)
	default:
		panic(fmt.Errorf("unknown prefix operator %s", node.Operator))
	}
}
