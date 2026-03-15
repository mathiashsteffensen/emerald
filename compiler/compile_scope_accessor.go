package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileScopeAccessor(node *ast.ScopeAccessor) {
	c.Compile(node.Left)

	c.emit(bytecode.OpScopedConstantGet, node.Token, c.addConstant(core.NewSymbol(node.Method.Value)))
}
