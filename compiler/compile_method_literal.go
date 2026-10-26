package compiler

import (
	"emerald/bytecode"
	"emerald/heap"
	"emerald/object"
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"strings"
)

func (c *Compiler) compileMethodLiteral(node *ast.MethodLiteral) {
	name := node.Name.(ast.IdentifierExpression).Value
	block, _ := c.compileBlock(node.BlockLiteral, true, name)

	symbol := c.rt.NewSymbol(name)

	c.emit(bytecode.OpPushConstant, node.Token, c.addConstant(symbol))
	c.emit(bytecode.OpPushConstant, node.BlockLiteral.Token, c.addConstant(object.NewHeapObject(block)))
	c.emit(bytecode.OpDefineMethod, node.Token)
}

// compileBlock compiles a method body (named name) or a block literal (name ignored) into a Block
func (c *Compiler) compileBlock(node *ast.BlockLiteral, method bool, name string) (*object.Block, int) {
	if !method {
		name = blockName(c.currentScope().name)
	}

	c.enterScope(method)
	c.scopes[c.scopeIndex].block = !method
	c.scopes[c.scopeIndex].name = name

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

		if rescueBlock.ErrorVarName != nil {
			c.emit(bytecode.OpRescuedError, rescueBlock.Token)
			c.compileAssignmentTarget(&ast.AssignmentExpression{Token: rescueBlock.Token, Name: rescueBlock.ErrorVarName})
			c.emit(bytecode.OpPop, rescueBlock.Token)
		}
		c.compileStatementsWithReturnValue(rescueBlock.Body.Statements, rescueBlock.Token)
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

	var kwargNames []string
	for _, argument := range node.KeywordArguments {
		kwargNames = append(kwargNames, argument.Value)
	}

	block := object.NewBlock(bytecode, numLocals, numParams, kwargNames, method)
	block.Name = name
	block.ExceptionTable = exceptionTable
	for _, symbol := range freeSymbols {
		block.FreeBindings = append(block.FreeBindings, object.FreeBinding{Index: symbol.Index, Local: symbol.Scope == heap.LocalScope})
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

// blockName labels a block literal by its enclosing method, like Ruby: "block in foo", "block (2 levels) in foo"
func blockName(enclosing string) string {
	levels := 1
	owner := enclosing

	if strings.HasPrefix(enclosing, "block in ") {
		levels = 2
		owner = strings.TrimPrefix(enclosing, "block in ")
	} else if n, rest, ok := parseBlockLevels(enclosing); ok {
		levels = n + 1
		owner = rest
	}

	if levels == 1 {
		return "block in " + owner
	}

	return fmt.Sprintf("block (%d levels) in %s", levels, owner)
}

func parseBlockLevels(name string) (int, string, bool) {
	var levels int
	var owner string
	if _, err := fmt.Sscanf(name, "block (%d levels) in %s", &levels, &owner); err != nil {
		return 0, "", false
	}
	// Sscanf stops %s at whitespace; the owner is everything after " in "
	if idx := strings.Index(name, " in "); idx >= 0 {
		owner = name[idx+len(" in "):]
	}
	return levels, owner, true
}
