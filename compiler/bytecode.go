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
	OpPop
	OpSub
	OpMul
	OpDiv
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpLessThan
	OpTrue
	OpFalse
	OpNull
	OpArray
	OpHash
	OpMinus
	OpBang
	OpJump
	OpJumpNotTruthy
	OpGetGlobal
	OpSetGlobal
	OpGetLocal
	OpSetLocal
	OpReturn
	OpReturnValue
	OpDefineMethod
	OpSend
)

var definitions = map[Opcode]*Definition{
	OpPushConstant:  {"OpPushConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpLessThan:      {"OpLessThan", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpNull:          {"OpNull", []int{}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJump:          {"OpJump", []int{2}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpGetLocal:      {"OpGetLocal", []int{1}},
	OpSetLocal:      {"OpSetLocal", []int{1}},
	OpReturn:        {"OpReturn", []int{}},
	OpReturnValue:   {"OpReturnValue", []int{}},
	OpDefineMethod:  {"OpDefineMethod", []int{}},
	OpSend:          {"OpSend", []int{1}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}
