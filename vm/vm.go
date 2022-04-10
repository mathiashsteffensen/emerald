package vm

import (
	"emerald/compiler"
	"emerald/object"
	"fmt"
)

const StackSize = 2048

// VM is our virtual machine responsible for the fetch, decode, execute cycle
type VM struct {
	constants    []object.EmeraldValue
	instructions compiler.Instructions

	stack []object.EmeraldValue
	sp    int // Always points to the next value. Top of stack is stack[sp-1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.EmeraldValue, StackSize),
		sp:           0,
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := compiler.Opcode(vm.instructions[ip])
		switch op {
		case compiler.OpAdd:
			right := vm.pop()
			left := vm.pop()

			result := left.(*object.IntegerInstance).Value + right.(*object.IntegerInstance).Value

			err := vm.push(object.NewInteger(result))
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
