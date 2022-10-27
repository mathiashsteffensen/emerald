package compiler

import (
	"emerald/core"
	"emerald/heap"
	"emerald/object"
	"emerald/parser"
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
)

type EmittedInstruction struct {
	Opcode   Opcode
	Position int
}

type Compiler struct {
	instructions Instructions
	opCount      int
	symbolTable  *heap.SymbolTable
	scopes       []CompilationScope
	scopeIndex   int
}

type ConstructorOption func(c *Compiler)

func New(options ...ConstructorOption) *Compiler {
	mainScope := CompilationScope{
		instructions:        Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c := &Compiler{
		instructions: Instructions{},
		symbolTable:  heap.GlobalSymbolTable,
		scopes:       []CompilationScope{mainScope},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func init() {
	core.Compile = func(fileName string, content string) []byte {
		l := lexer.New(lexer.NewInput(fileName, content))
		p := parser.New(l)
		ast := p.ParseAST()

		if len(p.Errors()) != 0 {
			core.Raise(core.NewException(fmt.Sprintf("failed to parse source file %s\n\n%s", fileName, p.Errors()[0])))
		}

		c := New()
		c.Compile(ast)

		instructions := append(c.Bytecode().Instructions, byte(OpReturn))

		return instructions
	}
}

func (c *Compiler) Compile(node ast.Node) {
	switch node := node.(type) {
	case *ast.AST:
		for _, s := range node.Statements {
			c.Compile(s)
		}
	case *ast.ExpressionStatement:
		c.Compile(node.Expression)
		c.emit(OpPop)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			c.Compile(s)
		}
	case *ast.ReturnStatement:
		c.Compile(node.ReturnValue)
		c.emit(OpReturnValue)
	case *ast.PrefixExpression:
		c.compilePrefixExpression(node)
	case *ast.AssignmentExpression:
		c.compileAssignment(node)
	case *ast.Self:
		c.emit(OpSelf)
	case ast.Yield:
		c.compileYield(node)
	case ast.IdentifierExpression:
		c.compileIdentifierExpression(node)
	case *ast.InstanceVariable:
		c.compileIdentifierExpression(node)
	case *ast.GlobalVariable:
		c.compileIdentifierExpression(node)
	case ast.CallExpression:
		c.emit(OpSelf) // Method calls without a receiver has an implicit self receiver
		c.compileCallExpression(node)
	case *ast.MethodCall:
		c.compileMethodCall(node)
	case *ast.ScopeAccessor:
		c.compileScopeAccessor(node)
	case *ast.InfixExpression:
		c.compileInfixExpression(node)
	case *ast.IfExpression:
		c.compileIfExpression(node)
	case *ast.CaseExpression:
		c.compileCaseExpression(node)

		if c.lastInstructionIs(OpPop) {
			c.removeLastPop()
		}
	case *ast.WhileExpression:
		c.compileWhileExpression(node)
	case *ast.IntegerLiteral:
		integer := core.NewInteger(node.Value)
		c.emit(OpPushConstant, c.addConstant(integer))
	case *ast.FloatLiteral:
		float := core.NewFloat(node.Value)
		c.emit(OpPushConstant, c.addConstant(float))
	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(OpTrue)
		} else {
			c.emit(OpFalse)
		}
	case *ast.NullExpression:
		c.emit(OpNull)
	case *ast.StringLiteral:
		c.compileStringLiteral(node)
	case *ast.StringTemplate:
		c.compileStringTemplate(node)
	case *ast.SymbolLiteral:
		sym := core.NewSymbol(node.Value)
		c.emit(OpPushConstant, c.addConstant(sym))
	case *ast.RegexpLiteral:
		c.compileRegexpLiteral(node)
	case *ast.ArrayLiteral:
		c.compileArrayLiteral(node)
	case *ast.HashLiteral:
		c.compileHashLiteral(node)
	case *ast.MethodLiteral:
		c.compileMethodLiteral(node)
	case *ast.ClassLiteral:
		c.compileClassLiteral(node)
	case *ast.StaticClassLiteral:
		c.compileStaticClassLiteral(node)
	case *ast.ModuleLiteral:
		c.compileModuleLiteral(node)
	}
}

func (c *Compiler) compileStatementsWithReturnValue(statements []ast.Statement) {
	if len(statements) == 0 {
		c.emit(OpNull)
	} else {
		for _, s := range statements {
			c.Compile(s)
		}
	}
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
	}
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = heap.NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) lastInstructionIs(op Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction
	old := c.currentInstructions()

	c.scopes[c.scopeIndex].instructions = old[:last.Position]
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) replaceLastInstructionWith(op Opcode) {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, Make(op))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = op
}

func (c *Compiler) changeOperand(opPos int, operands ...int) {
	op := Opcode(c.currentInstructions()[opPos])
	newInstruction := Make(op, operands...)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) emit(op Opcode, operands ...int) int {
	ins := Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	c.opCount += 1

	return pos
}

func (c *Compiler) emitConstantGet(name string) {
	symbol := core.NewSymbol(name)

	c.emit(OpConstantGet, c.addConstant(symbol))
}

func (c *Compiler) emitSymbol(symbol heap.Symbol) {
	switch symbol.Scope {
	case heap.GlobalScope:
		c.emit(OpGetGlobal, symbol.Index)
	case heap.LocalScope:
		c.emit(OpGetLocal, symbol.Index)
	case heap.FreeScope:
		c.emit(OpGetFree, symbol.Index)
	}
}

// returns the instructions for the current CompilationScope
func (c *Compiler) currentInstructions() Instructions {
	return c.scopes[c.scopeIndex].instructions
}

// addInstruction adds instructions to the instruction stack and returns its location
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

// addConstant adds a constant to the constant stack and returns its location
func (c *Compiler) addConstant(obj object.EmeraldValue) int {
	return heap.AddConstant(obj)
}
