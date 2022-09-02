package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileCallExpression(node ast.CallExpression) error {
	method := core.NewSymbol(node.Method.Value)

	c.emit(OpPushConstant, c.addConstant(method))

	if node.Block != nil {
		block, freeSymbolCount, err := c.compileBlock(node.Block)
		if err != nil {
			return err
		}

		c.emit(OpCloseBlock, c.addConstant(block), freeSymbolCount)
	} else {
		c.emit(OpNull)
	}

	for _, argument := range node.Arguments {
		err := c.Compile(argument)
		if err != nil {
			return err
		}
	}

	c.emit(OpSend, len(node.Arguments))

	return nil
}
