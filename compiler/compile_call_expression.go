package compiler

import (
	"emerald/core"
	"emerald/parser/ast"
)

func (c *Compiler) compileCallExpression(node ast.CallExpression) {
	method := core.NewSymbol(node.Method.Value)

	c.emit(OpPushConstant, c.addConstant(method))

	if node.Block != nil {
		block, freeSymbolCount := c.compileBlock(node.Block)

		c.emit(OpCloseBlock, c.addConstant(block), freeSymbolCount)
	} else {
		c.emit(OpNull)
	}

	for _, argument := range node.Arguments {
		c.Compile(argument)
	}

	hasKwargsOperand := 0 // 0 signals to VM that OpSend did not receive kwargs
	numKwargs := len(node.KeywordArguments)
	if numKwargs != 0 {
		for _, el := range node.KeywordArguments {
			c.Compile(el.Key)
			c.Compile(el.Value)
		}
		c.emit(OpHash, numKwargs*2)
		hasKwargsOperand = 1
	}

	c.emit(OpSend, len(node.Arguments)+numKwargs, hasKwargsOperand)
}
