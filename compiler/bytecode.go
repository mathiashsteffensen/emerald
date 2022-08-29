package compiler

import (
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
	}
)

func (bytecode Bytecode) String() string { return bytecode.Instructions.String() }

const (
	_ Opcode = iota

	// OpPushConstant pushes a constant from the constant pool onto the stack
	OpPushConstant
	OpPop

	// Infix operators
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMatch
	OpSpaceship
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanOrEq
	OpLessThan
	OpLessThanOrEq

	// Prefix operators
	OpMinus
	OpBang

	// OpSelf pushes the current value of self onto the stack
	OpSelf
	// OpTrue pushes core.TRUE onto the stack
	OpTrue
	// OpFalse pushes core.FALSE onto the stack
	OpFalse
	// OpNull pushes core.NULL onto the stack
	OpNull

	OpArray
	OpHash
	OpJump
	OpJumpTruthy
	OpJumpNotTruthy

	// OpConstantSet sets a constant in the current namespace.
	// By calling .NamespaceDefinitionSet on the current execution context
	OpConstantSet
	// OpConstantGet gets a constant from the current namespace.
	// By calling .NamespaceDefinitionGet on the current execution context
	OpConstantGet
	// OpConstantGetOrSet either gets a constant from the current namespace, if it is defined, otherwise defines it
	OpConstantGetOrSet
	// OpScopedConstantGet looks up constant in the element on top of the stack
	OpScopedConstantGet

	// OpGetGlobal resolves a global variable reference
	OpGetGlobal
	// OpSetGlobal creates a global variable reference
	OpSetGlobal

	// OpGetLocal resolves a local variable reference
	OpGetLocal
	// OpSetLocal creates a local variable reference
	OpSetLocal

	// OpGetFree like it's local & global variants resolves a variable reference, but from a blocks free variable pool
	OpGetFree

	OpReturn
	OpReturnValue
	OpDefineMethod

	// OpSend Invokes a method on the current execution context self.
	// Takes an operand that references the number of arguments passed.
	// Pops that number of arguments from the stack,
	// the next object on the stack is the symbol representing the name of the method to invoke
	OpSend

	// OpSetExecutionTarget takes the value at the top of the stack
	// and sets it as execution target
	// While replacing it in the stack with the previous target
	// so the previous target is the top of the stack
	OpSetExecutionTarget

	// OpResetExecutionTarget takes the value at the second-most top of the stack
	// and sets it as execution target
	OpResetExecutionTarget

	// OpOpenClass takes the value at the top of the stack
	// and sets it as definition context while setting its SingletonClass as execution context.
	OpOpenClass

	// OpCloseClass resets the execution context to its outer context
	OpCloseClass

	// Modifies whether the execution context is static.
	OpStaticTrue
	OpStaticFalse

	// OpCloseBlock creates a closure over a Block by fetching its free variables from the stack
	// and adding them to object.Block#FreeVariables
	// Has 2 operands, constant index of the block & number of free variables
	// NOTE: second operand is only 1 byte, so there is a hard limit on 256 free variables
	OpCloseBlock

	// Setter and getter operators for instance variables,
	// both take a single operand that's a reference to a constant in the constant pool
	// which is a symbol referencing the name.
	// Set operation sets value to top of the stack
	OpInstanceVarSet
	OpInstanceVarGet
)

var definitions = map[Opcode]*Definition{
	OpPushConstant:         {"OpPushConstant", []int{2}},
	OpAdd:                  {"OpAdd", []int{}},
	OpPop:                  {"OpPop", []int{}},
	OpSub:                  {"OpSub", []int{}},
	OpMul:                  {"OpMul", []int{}},
	OpDiv:                  {"OpDiv", []int{}},
	OpMatch:                {"OpMatch", []int{}},
	OpSpaceship:            {"OpSpaceship", []int{}},
	OpEqual:                {"OpEqual", []int{}},
	OpNotEqual:             {"OpNotEqual", []int{}},
	OpGreaterThan:          {"OpGreaterThan", []int{}},
	OpGreaterThanOrEq:      {"OpGreaterThanOrEq", []int{}},
	OpLessThan:             {"OpLessThan", []int{}},
	OpLessThanOrEq:         {"OpLessThanOrEq", []int{}},
	OpSelf:                 {"OpSelf", []int{}},
	OpTrue:                 {"OpTrue", []int{}},
	OpFalse:                {"OpFalse", []int{}},
	OpNull:                 {"OpNull", []int{}},
	OpArray:                {"OpArray", []int{2}},
	OpHash:                 {"OpHash", []int{2}},
	OpMinus:                {"OpMinus", []int{}},
	OpBang:                 {"OpBang", []int{}},
	OpJump:                 {"OpJump", []int{2}},
	OpJumpTruthy:           {"OpJumpTruthy", []int{2}},
	OpJumpNotTruthy:        {"OpJumpNotTruthy", []int{2}},
	OpGetGlobal:            {"OpGetGlobal", []int{2}},
	OpSetGlobal:            {"OpSetGlobal", []int{2}},
	OpGetLocal:             {"OpGetLocal", []int{1}},
	OpSetLocal:             {"OpSetLocal", []int{1}},
	OpGetFree:              {"OpGetFree", []int{1}},
	OpReturn:               {"OpReturn", []int{}},
	OpReturnValue:          {"OpReturnValue", []int{}},
	OpDefineMethod:         {"OpDefineMethod", []int{}},
	OpSend:                 {"OpSend", []int{1}},
	OpOpenClass:            {"OpOpenClass", []int{}},
	OpCloseClass:           {"OpCloseClass", []int{}},
	OpSetExecutionTarget:   {"OpSetExecutionTarget", []int{}},
	OpResetExecutionTarget: {"OpResetExecutionTarget", []int{}},
	OpStaticTrue:           {"OpStaticTrue", []int{}},
	OpStaticFalse:          {"OpStaticFalse", []int{}},
	OpCloseBlock:           {"OpCloseBlock", []int{2, 1}},
	OpInstanceVarSet:       {"OpInstanceVarSet", []int{2}},
	OpInstanceVarGet:       {"OpInstanceVarGet", []int{2}},
	OpConstantSet:          {"OpConstantSet", []int{2}},
	OpConstantGet:          {"OpConstantGet", []int{2}},
	OpConstantGetOrSet:     {"OpConstantGetOrSet", []int{2, 2}},
	OpScopedConstantGet:    {"OpScopedConstantGet", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}
