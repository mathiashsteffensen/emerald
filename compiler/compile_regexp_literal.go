package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileRegexpLiteral(node *ast.RegexpLiteral) {
	regexp := c.rt.NewRegexp(node.Value)

	c.emit(bytecode.OpPushConstant, node.Token, c.addConstant(regexp))
}
