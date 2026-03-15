package compiler

import (
	"emerald/bytecode"
	ast "emerald/parser/ast"
)

func (c *Compiler) compileHashLiteral(node *ast.HashLiteral) {
	for _, el := range node.Values {
		c.Compile(el.Key)
		c.Compile(el.Value)
	}
	c.emit(bytecode.OpHash, node.Token, len(node.Values)*2)
}
