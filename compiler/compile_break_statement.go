package compiler

import (
	"emerald/bytecode"
	"emerald/parser/ast"
)

func (c *Compiler) compileBreakStatement(node *ast.BreakStatement) {
	scope := &c.scopes[c.scopeIndex]

	// Inside a while loop, `break` jumps past the loop, its value becoming the loop's value
	if len(scope.loops) > 0 {
		c.Compile(node.Value)
		innermost := &scope.loops[len(scope.loops)-1]
		innermost.breakJumps = append(innermost.breakJumps, c.emit(bytecode.OpJump, node.Token, 9999))
		return
	}

	// Inside a block, `break` returns its value from the method call the block was passed to
	if !scope.block {
		c.abort("Invalid break")
		return
	}

	c.Compile(node.Value)
	c.emit(bytecode.OpBreak, node.Token)
}
