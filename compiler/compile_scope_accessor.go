package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileScopeAccessor(node *ast.ScopeAccessor) {
	c.Compile(node.Left)

	c.emit(OpScopedConstantGet, c.addConstant(core.NewSymbol(node.Method.Value)))
}
