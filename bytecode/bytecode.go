package bytecode

import (
	"emerald/parser/lexer"
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

		// DebugTokens maps instruction offsets to their originating tokens for debugging purposes
		DebugTokens map[int]lexer.Token
		Lexer       *lexer.Lexer
	}
)

func (bytecode Bytecode) String() string { return bytecode.Instructions.String() }

// InstructionSnapshot returns a snapshot of the code for the instruction at the given offset
func (bytecode Bytecode) InstructionSnapshot(offset int) string {
	return bytecode.Lexer.Snapshot(bytecode.DebugTokens[offset])
}

const (
	_ Opcode = iota

	// OpPushConstant pushes a constant from the constant pool onto the stack
	OpPushConstant
	OpPop
	// OpDupN duplicates the top N values, preserving their order.
	OpDupN
	// OpDropN removes N values below the top value.
	OpDropN
	// OpSwap exchanges the top two values.
	OpSwap

	// Infix operators
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMatch
	OpSpaceship
	OpEqual
	OpCaseEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanOrEq
	OpLessThan
	OpLessThanOrEq
	OpBinShiftLeft

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

	// OpStringJoin Takes an argument with the number of objects to fetch from the stack
	// It then calls #to_s on them & joins the strings
	// (Used by string templating)
	OpStringJoin

	// OpYield takes an argument that is the number of arguments to yield to the currently given block, if any
	OpYield

	// OpArray takes an argument that is the number of objects to take from the stack and add to a newly allocated array
	OpArray
	OpHash
	OpJump
	OpJumpTruthy
	OpJumpNotTruthy
	OpCheckCaseEqual

	// OpConstantSet sets a constant in the current namespace.
	// By calling .NamespaceDefinitionSet on self
	OpConstantSet
	// OpConstantGet gets a constant from the current namespace.
	// By calling .NamespaceDefinitionGet on self
	OpConstantGet
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
	OpSetFree

	OpReturn
	OpReturnValue
	OpRescuedError
	OpDefineMethod

	// OpSend invokes a method.
	// Takes operands for the method name constant, number of arguments, and whether keyword arguments were passed.
	OpSend
	// OpSendAssign invokes a setter but evaluates to its last argument.
	OpSendAssign

	// OpOpenClass takes an argument with a constant index pointing to a symbol name of the class to be set as self,
	// If no class exists with the specified name, it creates one.
	// the value at the top of the stack is the parent class.
	OpOpenClass

	// OpOpenModule takes an argument with a constant index pointing to a symbol name of the module to be set as self
	// If no module exists with the specified name, it creates one.
	OpOpenModule

	// OpUnwrapContext sets the execution context to its outer context
	OpUnwrapContext

	// Modifies whether the execution context is static.
	OpStaticTrue
	OpStaticFalse

	// OpCloseBlock creates a closure using the Block's free-variable bindings.
	// Captured locals share heap cells with their defining frame.
	// Has 2 operands, constant index of the block & number of free variables
	// NOTE: second operand is only 1 byte, so there is a hard limit on 256 free variables
	OpCloseBlock

	// Setter and getter operators for instance variables,
	// both take a single operand that's a reference to a constant in the constant pool.
	// That constant must be a symbol referencing the name.
	// Set operation sets value and pushes it to top of the stack
	OpInstanceVarSet
	OpInstanceVarGet
)

var definitions = map[Opcode]*Definition{
	OpPushConstant:      {"OpPushConstant", []int{2}},
	OpAdd:               {"OpAdd", []int{}},
	OpPop:               {"OpPop", []int{}},
	OpDupN:              {"OpDupN", []int{1}},
	OpDropN:             {"OpDropN", []int{1}},
	OpSwap:              {"OpSwap", []int{}},
	OpSub:               {"OpSub", []int{}},
	OpMul:               {"OpMul", []int{}},
	OpDiv:               {"OpDiv", []int{}},
	OpMatch:             {"OpMatch", []int{}},
	OpSpaceship:         {"OpSpaceship", []int{}},
	OpEqual:             {"OpEqual", []int{}},
	OpCaseEqual:         {"OpCaseEqual", []int{}},
	OpNotEqual:          {"OpNotEqual", []int{}},
	OpGreaterThan:       {"OpGreaterThan", []int{}},
	OpGreaterThanOrEq:   {"OpGreaterThanOrEq", []int{}},
	OpLessThan:          {"OpLessThan", []int{}},
	OpLessThanOrEq:      {"OpLessThanOrEq", []int{}},
	OpSelf:              {"OpSelf", []int{}},
	OpTrue:              {"OpTrue", []int{}},
	OpFalse:             {"OpFalse", []int{}},
	OpNull:              {"OpNull", []int{}},
	OpYield:             {"OpYield", []int{1}},
	OpArray:             {"OpArray", []int{2}},
	OpHash:              {"OpHash", []int{2}},
	OpMinus:             {"OpMinus", []int{}},
	OpBinShiftLeft:      {"OpBinShiftLeft", []int{}},
	OpBang:              {"OpBang", []int{}},
	OpJump:              {"OpJump", []int{2}},
	OpJumpTruthy:        {"OpJumpTruthy", []int{2}},
	OpJumpNotTruthy:     {"OpJumpNotTruthy", []int{2}},
	OpCheckCaseEqual:    {"OpCheckCaseEqual", []int{1, 2}},
	OpGetGlobal:         {"OpGetGlobal", []int{2}},
	OpSetGlobal:         {"OpSetGlobal", []int{2}},
	OpGetLocal:          {"OpGetLocal", []int{1}},
	OpSetLocal:          {"OpSetLocal", []int{1}},
	OpGetFree:           {"OpGetFree", []int{1}},
	OpSetFree:           {"OpSetFree", []int{1}},
	OpReturn:            {"OpReturn", []int{}},
	OpReturnValue:       {"OpReturnValue", []int{}},
	OpRescuedError:      {"OpRescuedError", []int{}},
	OpDefineMethod:      {"OpDefineMethod", []int{}},
	OpSend:              {"OpSend", []int{2, 1, 1}},
	OpSendAssign:        {"OpSendAssign", []int{2, 1, 1}},
	OpOpenClass:         {"OpOpenClass", []int{2}},
	OpOpenModule:        {"OpOpenModule", []int{2}},
	OpUnwrapContext:     {"OpUnwrapContext", []int{}},
	OpStaticTrue:        {"OpStaticTrue", []int{}},
	OpStaticFalse:       {"OpStaticFalse", []int{}},
	OpCloseBlock:        {"OpCloseBlock", []int{2, 1}},
	OpInstanceVarSet:    {"OpInstanceVarSet", []int{2}},
	OpInstanceVarGet:    {"OpInstanceVarGet", []int{2}},
	OpConstantSet:       {"OpConstantSet", []int{2}},
	OpConstantGet:       {"OpConstantGet", []int{2}},
	OpScopedConstantGet: {"OpScopedConstantGet", []int{2}},
	OpStringJoin:        {"OpStringJoin", []int{1}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}
