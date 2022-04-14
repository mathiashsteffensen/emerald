package vm

import (
	"emerald/compiler"
	"emerald/object"
	"fmt"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
	MaxFrames   = 1024
)

// VM is our virtual machine responsible for the fetch, decode, execute cycle
type VM struct {
	ec          object.ExecutionContext
	constants   []object.EmeraldValue
	stack       []object.EmeraldValue
	sp          int // Always points to the next value. Top of stack is stack[sp-1]
	globals     []object.EmeraldValue
	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainBlock := &object.Block{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainBlock, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		ec:          object.ExecutionContext{Target: object.Object, IsStatic: true},
		constants:   bytecode.Constants,
		stack:       make([]object.EmeraldValue, StackSize),
		sp:          0,
		globals:     make([]object.EmeraldValue, GlobalsSize),
		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.EmeraldValue) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) Run() error {
	var (
		ip  int
		ins compiler.Instructions
		op  compiler.Opcode
		err error
	)

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = compiler.Opcode(ins[ip])

		switch op {
		case compiler.OpPop:
			vm.pop()
		case compiler.OpTrue:
			err = vm.push(object.TRUE)
		case compiler.OpFalse:
			err = vm.push(object.FALSE)
		case compiler.OpNull:
			err = vm.push(object.NULL)
		case compiler.OpPushConstant:
			constIndex := compiler.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err = vm.push(vm.constants[constIndex])
		case compiler.OpJump:
			pos := int(compiler.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case compiler.OpJumpNotTruthy:
			pos := int(compiler.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case compiler.OpGetGlobal:
			globalIndex := compiler.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			err = vm.push(vm.globals[globalIndex])
		case compiler.OpSetGlobal:
			globalIndex := compiler.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()
		case compiler.OpGetLocal:
			localIndex := compiler.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			err = vm.push(vm.stack[frame.basePointer+int(localIndex)])
		case compiler.OpSetLocal:
			localIndex := compiler.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
		case compiler.OpArray:
			numElements := int(compiler.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err = vm.push(array)
		case compiler.OpBang:
			err = vm.executeBangOperator()
		case compiler.OpMinus:
			err = vm.executeMinusOperator()
		case compiler.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err = vm.push(object.NULL)
		case compiler.OpReturnValue:
			returnValue := vm.pop()
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err = vm.push(returnValue)
		case compiler.OpDefineMethod:
			name := vm.pop().(*object.SymbolInstance)
			block := vm.pop().(*object.Block)

			result := vm.ec.Target.DefineMethod(vm.ec.IsStatic, block, name)

			err = vm.push(result)
		case compiler.OpSend:
			numArgs := compiler.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			err = vm.callFunction(int(numArgs))
		default:
			if isInfixOperator(op) {
				left, operator, right := vm.decodeInfixOperator(op)

				err = vm.executeInfixOperator(operator, left, right)
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (vm *VM) callFunction(numArgs int) error {
	name := vm.stack[vm.sp-1-numArgs].(*object.SymbolInstance)

	target := vm.ec.Target
	method, errVal := target.ExtractMethod(name.Value, target, target)
	if errVal != nil {
		return vm.push(errVal)
	} else {
		switch method := method.(type) {
		case *object.Block:
			frame := NewFrame(method, vm.sp-numArgs)
			vm.pushFrame(frame)
			vm.sp = frame.basePointer + method.NumLocals
		case *object.WrappedBuiltInMethod:
			method.Method(target, nil, nil, vm.stack[vm.sp-numArgs:vm.sp]...)
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
