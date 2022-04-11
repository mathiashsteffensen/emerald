package compiler

import (
	"emerald/ast"
	"emerald/object"
	"fmt"
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
}

func New() *Compiler {
	return &Compiler{
		instructions:        Instructions{},
		constants:           []object.EmeraldValue{},
		symbolTable:         NewSymbolTable(),
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
}

func NewWithState(s *SymbolTable, constants []object.EmeraldValue) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
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
	case *ast.PrefixExpression:
		err := c.compilePrefixExpression(node)
		if err != nil {
			return err
		}
	case *ast.AssignmentExpression:
		err := c.compileGlobalAssignment(node)
		if err != nil {
			return err
		}
	case *ast.IdentifierExpression:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.emit(OpGetGlobal, symbol.Index)
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
		integer := object.NewInteger(node.Value)
		c.emit(OpPushConstant, c.addConstant(integer))
	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(OpTrue)
		} else {
			c.emit(OpFalse)
		}
	case *ast.StringLiteral:
		str := object.NewString(node.Value)
		c.emit(OpPushConstant, c.addConstant(str))
	case *ast.ArrayLiteral:
		err := c.compileArrayLiteral(node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := Opcode(c.instructions[opPos])
	newInstruction := Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) emit(op Opcode, operands ...int) int {
	ins := Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

// addInstruction adds instructions to the instruction stack and returns its location
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	c.previousInstruction = previous
	c.lastInstruction = last
}

// addConstant adds a constant to the constant stack and returns its location
func (c *Compiler) addConstant(obj object.EmeraldValue) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}
