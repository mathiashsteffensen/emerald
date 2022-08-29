package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileRegexpLiteral(node *ast.RegexpLiteral) {
	regexp := core.NewRegexp(node.Value)

	c.emit(OpPushConstant, c.addConstant(regexp))
}
