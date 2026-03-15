package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileRegexpLiteral(node *ast.RegexpLiteral) {
	regexp := core.NewRegexp(node.Value)

	c.emit(bytecode.OpPushConstant, node.Token, c.addConstant(regexp))
}
