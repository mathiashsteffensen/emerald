package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileScopeAccessor(node *ast.ScopeAccessor) {
	c.Compile(node.Left)

	c.emit(bytecode.OpScopedConstantGet, node.Token, c.addConstant(c.rt.NewSymbol(node.Method.Value)))
}
