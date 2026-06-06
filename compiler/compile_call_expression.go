package compiler

import (
	"emerald/bytecode"
	"emerald/object"
	"emerald/parser/ast"
)

func (c *Compiler) compileCallExpression(node ast.CallExpression) {
	methodIndex := c.addConstant(c.rt.NewSymbol(node.Method.Value))

	if node.Block != nil {
		block, freeSymbolCount := c.compileBlock(node.Block, false)

		c.emit(bytecode.OpCloseBlock, node.Block.Token, c.addConstant(object.NewHeapObject(block)), freeSymbolCount)
	} else {
		c.emit(bytecode.OpNull, node.Token)
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
		c.emit(bytecode.OpHash, node.Token, numKwargs*2)
		hasKwargsOperand = 1
	}

	c.emit(bytecode.OpSend, node.Token, methodIndex, len(node.Arguments)+numKwargs, hasKwargsOperand)
}
