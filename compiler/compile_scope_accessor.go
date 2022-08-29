package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileScopeAccessor(node *ast.ScopeAccessor) {
	c.emit(OpConstantGet, c.addConstant(core.NewSymbol(node.Left.String())))
	c.emit(OpScopedConstantGet, c.addConstant(core.NewSymbol(node.Method.Value)))
}
