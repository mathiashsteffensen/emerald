package vm

import (
	"emerald/compiler"
	"emerald/object"
	"fmt"
)

const StackSize = 2048
const GlobalsSize = 65536

// VM is our virtual machine responsible for the fetch, decode, execute cycle
type VM struct {
	constants    []object.EmeraldValue
	instructions compiler.Instructions

	stack []object.EmeraldValue
	sp    int // Always points to the next value. Top of stack is stack[sp-1]

	globals []object.EmeraldValue
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.EmeraldValue, StackSize),
		sp:           0,
		globals:      make([]object.EmeraldValue, GlobalsSize),
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.EmeraldValue) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := compiler.Opcode(vm.instructions[ip])

		switch op {
		case compiler.OpPop:
			vm.pop()
		case compiler.OpTrue:
			err := vm.push(object.TRUE)
			if err != nil {
				return err
			}
		case compiler.OpFalse:
			err := vm.push(object.FALSE)
			if err != nil {
				return err
			}
		case compiler.OpNull:
			err := vm.push(object.NULL)
			if err != nil {
				return err
			}
		case compiler.OpPushConstant:
			constIndex := compiler.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case compiler.OpJump:
			pos := int(compiler.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case compiler.OpJumpNotTruthy:
			pos := int(compiler.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}
		case compiler.OpGetGlobal:
			globalIndex := compiler.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case compiler.OpSetGlobal:
			globalIndex := compiler.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIndex] = vm.pop()
		case compiler.OpArray:
			numElements := int(compiler.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}
		case compiler.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case compiler.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}
		default:
			if isInfixOperator(op) {
				left, operator, right := vm.decodeInfixOperator(op)

				err := vm.executeInfixOperator(operator, left, right)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// StackTop fetches the object at the top of the stack
func (vm *VM) StackTop() object.EmeraldValue {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.EmeraldValue {
	return vm.stack[vm.sp]
}

// push an obj on to the stack
func (vm *VM) push(obj object.EmeraldValue) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow: max stack size of %d exceeded", StackSize)
	}

	vm.stack[vm.sp] = obj
	vm.sp++

	return nil
}

// pop an obj from the top of the stack
func (vm *VM) pop() object.EmeraldValue {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) buildArray(startIndex, endIndex int) object.EmeraldValue {
	elements := make([]object.EmeraldValue, endIndex-startIndex)
	
	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.ArrayInstance{Value: elements}
}

func isTruthy(obj object.EmeraldValue) bool {
	switch obj {
	case object.FALSE, object.NULL:
		return false
	}

	return true
}
