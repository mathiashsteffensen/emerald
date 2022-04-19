package compiler

import (
	"emerald/ast"
	"emerald/core"
	"emerald/object"
)

type EmittedInstruction struct {
	Opcode   Opcode
	Position int
}

type Compiler struct {
	instructions        Instructions
	constants           []object.EmeraldValue
	symbolTable         *SymbolTable
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	scopes              []CompilationScope
	scopeIndex          int
}

type ConstructorOption func(c *Compiler)

func New(options ...ConstructorOption) *Compiler {
	mainScope := CompilationScope{
		instructions:        Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c := &Compiler{
		instructions:        Instructions{},
		constants:           []object.EmeraldValue{},
		symbolTable:         NewSymbolTable(),
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		scopes:              []CompilationScope{mainScope},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithState(s *SymbolTable, constants []object.EmeraldValue) ConstructorOption {
	return func(c *Compiler) {
		c.symbolTable = s
		c.constants = constants
	}
}

func WithBuiltIns() ConstructorOption {
	return func(c *Compiler) {
		for key, value := range object.Classes {
			symbol, ok := c.symbolTable.Resolve(key)
			if !ok {
				symbol = c.symbolTable.Define(key)

				c.emit(OpPushConstant, c.addConstant(value))
				c.emit(OpSetGlobal, symbol.Index)
				c.emit(OpPop)
			}
		}
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.AST:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}

		c.emit(OpPop)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(OpReturnValue)
	case *ast.PrefixExpression:
		err := c.compilePrefixExpression(node)
		if err != nil {
			return err
		}
	case *ast.AssignmentExpression:
		err := c.compileAssignment(node)
		if err != nil {
			return err
		}
	case *ast.IdentifierExpression:
		c.compileIdentifierExpression(node)
	case *ast.InstanceVariable:
		c.compileIdentifierExpression(node)
	case *ast.CallExpression:
		err := c.compileCallExpression(node)
		if err != nil {
			return err
		}
	case *ast.MethodCall:
		err := c.compileMethodCall(node)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
		err := c.compileInfixExpression(node)
		if err != nil {
			return err
		}
	case *ast.IfExpression:
		err := c.compileIfExpression(node)
		if err != nil {
			return err
		}
	case *ast.IntegerLiteral:
		integer := core.NewInteger(node.Value)
		c.emit(OpPushConstant, c.addConstant(integer))
	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(OpTrue)
		} else {
			c.emit(OpFalse)
		}
	case *ast.StringLiteral:
		str := core.NewString(node.Value)
		c.emit(OpPushConstant, c.addConstant(str))
	case *ast.SymbolLiteral:
		sym := core.NewSymbol(node.Value)
		c.emit(OpPushConstant, c.addConstant(sym))
	case *ast.ArrayLiteral:
		err := c.compileArrayLiteral(node)
		if err != nil {
			return err
		}
	case *ast.HashLiteral:
		err := c.compileHashLiteral(node)
		if err != nil {
			return err
		}
	case *ast.MethodLiteral:
		err := c.compileMethodLiteral(node)
		if err != nil {
			return err
		}
	case *ast.ClassLiteral:
		err := c.compileClassLiteral(node)
		if err != nil {
			return err
		}
	case *ast.StaticClassLiteral:
		err := c.compileStaticClassLiteral(node)
		if err != nil {
			return err
		}
	case *ast.ModuleLiteral:
		err := c.compileModuleLiteral(node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
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

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
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

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := Opcode(c.currentInstructions()[opPos])
	newInstruction := Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) emit(op Opcode, operands ...int) int {
	ins := Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) emitSymbol(symbol Symbol) {
	switch symbol.Scope {
	case GlobalScope:
		c.emit(OpGetGlobal, symbol.Index)
	case LocalScope:
		c.emit(OpGetLocal, symbol.Index)
	case FreeScope:
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
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}
