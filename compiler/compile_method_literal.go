package compiler

import (
	"emerald/ast"
	"emerald/object"
)

func (c *Compiler) compileMethodLiteral(node *ast.MethodLiteral) error {
	c.enterScope()

	for _, p := range node.Parameters {
		c.symbolTable.Define(p.(*ast.IdentifierExpression).Value)
	}

	err := c.Compile(node.Body)
	if err != nil {
		return err
	}

	if c.lastInstructionIs(OpPop) {
		c.replaceLastPopWithReturn()
	}

	// Everything in Ruby returns something
	// If function doesn't have a return value, return null
	if !c.lastInstructionIs(OpReturnValue) {
		c.emit(OpReturn)
	}

	numLocals := c.symbolTable.numDefinitions
	instructions := c.leaveScope()
	symbol := object.NewSymbol(node.Name.(*ast.IdentifierExpression).Value)
	block := object.NewBlock([]ast.Expression{}, instructions, numLocals)

	c.emit(OpPushConstant, c.addConstant(symbol))
	c.emit(OpPushConstant, c.addConstant(block))
	c.emit(OpDefineMethod)

	return nil
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, Make(OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = OpReturnValue
}
