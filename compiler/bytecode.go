package compiler

import (
	"emerald/object"
	"fmt"
)

type (
	// Opcode is just a byte, so it can be part of the instructions
	// We use an aliased type so the Go compiler can tell us if we accidentally try to use a value as an operator
	Opcode byte

	// Definition is a type for us to store detailed definitions or metadata for our Opcodes
	Definition struct {
		Name          string
		OperandWidths []int
	}

	Bytecode struct {
		Instructions Instructions
		Constants    []object.EmeraldValue
	}
)

const (
	OpPushConstant Opcode = iota
	OpAdd
)

var definitions = map[Opcode]*Definition{
	OpPushConstant: {"OpPushConstant", []int{2}},
	OpAdd:          {"OpAdd", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}
