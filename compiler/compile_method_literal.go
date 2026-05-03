package compiler

import (
	"emerald/bytecode"
	"emerald/object"
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
)

func (c *Compiler) compileMethodLiteral(node *ast.MethodLiteral) {
	block, _ := c.compileBlock(node.BlockLiteral, true)

	symbol := c.rt.NewSymbol(node.Name.(ast.IdentifierExpression).Value)

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
	startIP := len(c.currentInstructions())

	c.Compile(node.Body)

	c.ensureLastInstructionIsReturn(node.Token)

	endIP := len(c.currentInstructions())

	jumpsToPatch := []int{}

	if len(node.RescueBlocks) > 0 {
		jumpsToPatch = append(jumpsToPatch, c.emit(bytecode.OpJump, node.Token, 0))
	}

	exceptionTable := []object.ExceptionTableEntry{}

	for i, rescueBlock := range node.RescueBlocks {
		handlerIP := len(c.currentInstructions())

		c.Compile(rescueBlock.Body)
		c.ensureLastInstructionIsReturn(rescueBlock.Token)

		// We only need to emit a jump if there are more rescue blocks to follow,
		// otherwise we'll just fall through to the end of the method anyway.
		if i < len(node.RescueBlocks)-1 {
			jumpsToPatch = append(jumpsToPatch, c.emit(bytecode.OpJump, rescueBlock.Token, 0))
		}

		var errorClasses []string
		for _, errorClass := range rescueBlock.CaughtErrorClasses {
			errorClasses = append(errorClasses, errorClass.String(0))
		}

		exceptionTable = append(exceptionTable, object.ExceptionTableEntry{
			StartIP:            startIP,
			EndIP:              endIP,
			HandlerIP:          handlerIP,
			CaughtErrorClasses: errorClasses,
		})
	}

	for _, pos := range jumpsToPatch {
		c.changeOperand(pos, len(c.currentInstructions()))
	}

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
	block.ExceptionTable = exceptionTable

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
