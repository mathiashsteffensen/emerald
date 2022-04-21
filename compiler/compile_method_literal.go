package compiler

import (
	"emerald/ast"
	"emerald/core"
	"emerald/object"
)

func (c *Compiler) compileMethodLiteral(node *ast.MethodLiteral) error {
	block, _, err := c.compileBlock(node.BlockLiteral)
	if err != nil {
		return err
	}

	symbol := core.NewSymbol(node.Name.(*ast.IdentifierExpression).Value)

	c.emit(OpPushConstant, c.addConstant(symbol))
	c.emit(OpPushConstant, c.addConstant(block))
	c.emit(OpDefineMethod)

	return nil
}

func (c *Compiler) compileBlock(node *ast.BlockLiteral) (*object.Block, int, error) {
	c.enterScope()

	numParams := len(node.Parameters)
	for _, p := range node.Parameters {
		c.symbolTable.Define(p.(*ast.IdentifierExpression).Value)
	}

	err := c.Compile(node.Body)
	if err != nil {
		return nil, 0, err
	}

	if c.lastInstructionIs(OpPop) {
		c.replaceLastPopWithReturn()
	}

	// Everything in Ruby returns something
	// If function doesn't have a return value, return null
	if !c.lastInstructionIs(OpReturnValue) {
		c.emit(OpReturn)
	}

	freeSymbols := c.symbolTable.FreeSymbols
	numLocals := c.symbolTable.numDefinitions
	instructions := c.leaveScope()

	for _, s := range freeSymbols {
		c.emitSymbol(s)
	}

	return object.NewBlock(instructions, numLocals, numParams), len(freeSymbols), nil
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, Make(OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = OpReturnValue
}
