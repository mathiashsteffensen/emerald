package compiler

import (
	"emerald/ast"
	"emerald/core"
)

func (c *Compiler) compileRegexpLiteral(node *ast.RegexpLiteral) {
	regexp := core.NewRegexp(node.Value)

	c.emit(OpPushConstant, c.addConstant(regexp))
}
