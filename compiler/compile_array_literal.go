package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileArrayLiteral(node *ast.ArrayLiteral) {
	for _, val := range node.Value {
		c.Compile(val)
	}

	c.emit(bytecode.OpArray, node.Token, len(node.Value))
}
