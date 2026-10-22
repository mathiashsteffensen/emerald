package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileNextStatement(node *ast.NextStatement) {
	scope := c.currentScope()

	// Inside a while loop, `next` skips to the condition. Its value is discarded.
	if len(scope.loops) > 0 {
		c.Compile(node.Value)
		c.emit(bytecode.OpPop, node.Token)
		c.emit(bytecode.OpJump, node.Token, scope.loops[len(scope.loops)-1].conditionPosition)
		return
	}

	// Inside a block, `next` returns its value from the block
	if !scope.block {
		c.abort("Invalid next")
		return
	}

	c.Compile(node.Value)
	c.emit(bytecode.OpReturnValue, node.Token)
}
