package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileScopeAccessor(node *ast.ScopeAccessor) error {
	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	c.emit(OpScopedConstantGet, c.addConstant(core.NewSymbol(node.Method.Value)))

	return nil
}
