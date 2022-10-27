package compiler

import (
	"emerald/core"
	"emerald/object"
	ast "emerald/parser/ast"
)

func (c *Compiler) compileMethodLiteral(node *ast.MethodLiteral) {
	block, _ := c.compileBlock(node.BlockLiteral)

	symbol := core.NewSymbol(node.Name.(ast.IdentifierExpression).Value)

	c.emit(OpPushConstant, c.addConstant(symbol))
	c.emit(OpPushConstant, c.addConstant(block))
	c.emit(OpDefineMethod)
}

func (c *Compiler) compileBlock(node *ast.BlockLiteral) (*object.Block, int) {
	c.enterScope()

	numParams := len(node.Arguments)
	for _, p := range append(node.Arguments, node.KeywordArguments...) {
		c.symbolTable.Define(p.Value)
	}

	c.Compile(node.Body)

	if c.lastInstructionIs(OpPop) {
		c.replaceLastPopWithReturn()
	}

	// Everything in Ruby returns something
	// If function doesn't have a return value, return null
	if !c.lastInstructionIs(OpReturnValue) {
		c.emit(OpReturn)
	}

	freeSymbols := c.symbolTable.FreeSymbols
	numLocals := c.symbolTable.NumDefinitions
	instructions := c.leaveScope()

	for _, s := range freeSymbols {
		c.emitSymbol(s)
	}

	var kwargNames []string
	for _, argument := range node.KeywordArguments {
		kwargNames = append(kwargNames, argument.Value)
	}

	block := object.NewBlock(instructions, numLocals, numParams, kwargNames)

	for _, rescueBlock := range node.RescueBlocks {
		c.enterScope()

		c.Compile(rescueBlock.Body)

		if c.lastInstructionIs(OpPop) {
			c.replaceLastPopWithReturn()
		}

		// Everything in Ruby returns something
		// If function doesn't have a return value, return null
		if !c.lastInstructionIs(OpReturnValue) {
			c.emit(OpReturn)
		}

		var errorClasses []string
		instructions = c.leaveScope()

		for _, errorClass := range rescueBlock.CaughtErrorClasses {
			errorClasses = append(errorClasses, errorClass.String(0))
		}

		block.RescueBlocks = append(block.RescueBlocks, object.NewRescueBlock(instructions, errorClasses...))
	}

	return block, len(freeSymbols)
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, Make(OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = OpReturnValue
}
