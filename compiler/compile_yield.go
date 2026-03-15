package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileYield(node ast.Yield) {
	for _, argument := range node.Arguments {
		c.Compile(argument)
	}

	c.emit(bytecode.OpYield, node.Token, len(node.Arguments))
}
