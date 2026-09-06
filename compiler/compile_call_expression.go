package compiler

import (
	"emerald/bytecode"
	"emerald/object"
	"emerald/parser/ast"
)

func (c *Compiler) compileCallExpression(node ast.CallExpression) {
	methodIndex := c.addConstant(c.rt.NewSymbol(node.Method.Value))
	numArgs, hasKwargs := c.compileCallArguments(node)
	op := bytecode.OpSend
	if node.Assignment {
		op = bytecode.OpSendAssign
	}
	c.emit(op, node.Token, methodIndex, numArgs, hasKwargs)
}

func (c *Compiler) compileCallArguments(node ast.CallExpression) (int, int) {
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

	return len(node.Arguments) + numKwargs, hasKwargsOperand
}
