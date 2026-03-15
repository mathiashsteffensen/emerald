package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/object"
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
)

func (c *Compiler) compileMethodLiteral(node *ast.MethodLiteral) {
	block, _ := c.compileBlock(node.BlockLiteral, true)

	symbol := core.NewSymbol(node.Name.(ast.IdentifierExpression).Value)

	c.emit(bytecode.OpPushConstant, node.Token, c.addConstant(symbol))
	c.emit(bytecode.OpPushConstant, node.BlockLiteral.Token, c.addConstant(block))
	c.emit(bytecode.OpDefineMethod, node.Token)
}

func (c *Compiler) compileBlock(node *ast.BlockLiteral, enforceArity bool) (*object.Block, int) {
	c.enterScope()

	numParams := len(node.Arguments)
	for _, p := range append(node.Arguments, node.KeywordArguments...) {
		c.symbolTable.Define(p.Value)
	}

	c.Compile(node.Body)

	c.ensureLastInstructionIsReturn(node.Token)

	freeSymbols := c.symbolTable.FreeSymbols
	numLocals := c.symbolTable.NumDefinitions
	bytecode := c.leaveScope()

	for _, s := range freeSymbols {
		c.emitSymbol(s, node.Token)
	}

	var kwargNames []string
	for _, argument := range node.KeywordArguments {
		kwargNames = append(kwargNames, argument.Value)
	}

	block := object.NewBlock(bytecode, numLocals, numParams, kwargNames, enforceArity)

	for _, rescueBlock := range node.RescueBlocks {
		c.enterScope()

		c.Compile(rescueBlock.Body)

		c.ensureLastInstructionIsReturn(rescueBlock.Token)

		var errorClasses []string
		bytecode = c.leaveScope()

		for _, errorClass := range rescueBlock.CaughtErrorClasses {
			errorClasses = append(errorClasses, errorClass.String(0))
		}

		block.RescueBlocks = append(block.RescueBlocks, object.NewRescueBlock(bytecode, errorClasses...))
	}

	return block, len(freeSymbols)
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, bytecode.Make(bytecode.OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = bytecode.OpReturnValue
}

func (c *Compiler) ensureLastInstructionIsReturn(token lexer.Token) {
	// These 2 if statements ensure everything in the Ruby implementation returns something

	// If last instruction is a pop (an expression that leaves a value on the stack),
	// replace it with return
	if c.lastInstructionIs(bytecode.OpPop) {
		c.replaceLastPopWithReturn()
	}

	// If block doesn't have a return value, return null
	if !c.lastInstructionIs(bytecode.OpReturnValue) {
		c.emit(bytecode.OpReturn, token)
	}
}
