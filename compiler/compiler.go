package compiler

import (
	"emerald/ast"
	"emerald/object"
	"fmt"
)

type Compiler struct {
	instructions Instructions
	constants    []object.EmeraldValue
}

func New() *Compiler {
	return &Compiler{
		instructions: Instructions{},
		constants:    []object.EmeraldValue{},
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
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(OpAdd)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := object.NewInteger(node.Value)
		c.emit(OpPushConstant, c.addConstant(integer))
	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) emit(op Opcode, operands ...int) int {
	ins := Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

// addInstruction adds instructions to the instruction stack and returns its location
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

// addConstant adds a constant to the constant stack and returns its location
func (c *Compiler) addConstant(obj object.EmeraldValue) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}
