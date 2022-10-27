package compiler

import (
	"emerald/parser/ast"
)

func (c *Compiler) compileArrayLiteral(node *ast.ArrayLiteral) {
	for _, val := range node.Value {
		c.Compile(val)
	}

	c.emit(OpArray, len(node.Value))
}
