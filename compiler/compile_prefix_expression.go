package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
	"fmt"
)

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) {
	c.Compile(node.Right)

	switch node.Operator {
	case "!":
		c.emit(bytecode.OpBang, node.Token)
	case "-":
		c.emit(bytecode.OpMinus, node.Token)
	default:
		panic(fmt.Errorf("unknown prefix operator %s", node.Operator))
	}
}
