package compiler

import "emerald/parser/ast"

func (c *Compiler) compileYield(node ast.Yield) {
	for _, argument := range node.Arguments {
		c.Compile(argument)
	}

	c.emit(OpYield, len(node.Arguments))
}
