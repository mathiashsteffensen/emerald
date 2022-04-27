package vm

import (
	"emerald/compiler"
	"emerald/core"
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
	ctx         *object.Context
	constants   []object.EmeraldValue
	stack       []object.EmeraldValue
	sp          int // Always points to the next value. Top of stack is stack[sp-1]
	globals     []object.EmeraldValue
	frames      []*Frame
	framesIndex int
}

type ConstructorOption func(vm *VM)

func New(bytecode *compiler.Bytecode, options ...ConstructorOption) *VM {
	mainBlock := &object.ClosedBlock{Block: &object.Block{Instructions: bytecode.Instructions}}
	mainFrame := NewFrame(mainBlock, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	vm := &VM{
		ctx:         &object.Context{DefinitionTarget: core.Object, ExecutionTarget: core.MainObject},
		constants:   bytecode.Constants,
		stack:       make([]object.EmeraldValue, StackSize),
		sp:          0,
		globals:     make([]object.EmeraldValue, GlobalsSize),
		frames:      frames,
		framesIndex: 1,
	}

	for _, option := range options {
		option(vm)
	}

	return vm
}

func WithGlobalsStore(s []object.EmeraldValue) ConstructorOption {
	return func(vm *VM) {
		vm.globals = s
	}
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

		ip, ins, op = vm.fetch()

		err = vm.execute(ip, ins, op)
		if err != nil {
			return err
		}
	}

	return err
}

func (vm *VM) fetch() (int, compiler.Instructions, compiler.Opcode) {
	ip := vm.currentFrame().ip
	ins := vm.currentFrame().Instructions()
	return ip, ins, compiler.Opcode(ins[ip])
}

func (vm *VM) execute(ip int, ins compiler.Instructions, op compiler.Opcode) error {
	switch op {
	case compiler.OpPop:
		vm.pop()
	case compiler.OpTrue:
		err := vm.push(core.TRUE)
		if err != nil {
			return err
		}
	case compiler.OpFalse:
		err := vm.push(core.FALSE)
		if err != nil {
			return err
		}
	case compiler.OpNull:
		err := vm.push(core.NULL)
		if err != nil {
			return err
		}
	case compiler.OpPushConstant:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		err := vm.push(vm.constants[constIndex])
		if err != nil {
			return err
		}
	case compiler.OpJump:
		pos := int(compiler.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip = pos - 1
	case compiler.OpJumpNotTruthy:
		vm.conditionalJump(!isTruthy(vm.StackTop()), ins, ip)
	case compiler.OpJumpTruthy:
		vm.conditionalJump(isTruthy(vm.StackTop()), ins, ip)
	case compiler.OpGetGlobal:
		globalIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2
		err := vm.push(vm.globals[globalIndex])
		if err != nil {
			return err
		}
	case compiler.OpSetGlobal:
		globalIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2
		vm.globals[globalIndex] = vm.StackTop()
	case compiler.OpGetLocal:
		localIndex := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1
		frame := vm.currentFrame()
		err := vm.push(vm.stack[frame.basePointer+int(localIndex)])
		if err != nil {
			return err
		}
	case compiler.OpSetLocal:
		localIndex := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1
		frame := vm.currentFrame()
		vm.stack[frame.basePointer+int(localIndex)] = vm.StackTop()
	case compiler.OpGetFree:
		freeIndex := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1

		err := vm.push(vm.currentFrame().block.FreeVariables[freeIndex])
		if err != nil {
			return err
		}
	case compiler.OpInstanceVarGet:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		name := vm.constants[constIndex]
		target := vm.ctx.ExecutionTarget

		val := target.InstanceVariableGet(name.(*core.SymbolInstance).Value, target, target)

		if val == nil {
			val = core.NULL
		}

		err := vm.push(val)
		if err != nil {
			return err
		}
	case compiler.OpInstanceVarSet:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		name := vm.constants[constIndex]
		val := vm.StackTop()
		target := vm.ctx.ExecutionTarget

		target.InstanceVariableSet(name.(*core.SymbolInstance).Value, val)
	case compiler.OpArray:
		numElements := int(compiler.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip += 2

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
	case compiler.OpReturn:
		frame := vm.popFrame()
		vm.sp = frame.basePointer - 2

		err := vm.push(core.NULL)
		if err != nil {
			return err
		}
	case compiler.OpReturnValue:
		returnValue := vm.pop()

		frame := vm.popFrame()
		vm.sp = frame.basePointer - 2

		err := vm.push(returnValue)
		if err != nil {
			return err
		}
	case compiler.OpDefineMethod:
		block := vm.pop().(*object.Block)
		name := vm.stack[vm.sp-1].(*core.SymbolInstance)

		vm.ctx.DefinitionTarget.DefineMethod(object.NewClosedBlock(block, []object.EmeraldValue{}), name)
	case compiler.OpSend:
		numArgs := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1
		err := vm.callFunction(int(numArgs))
		if err != nil {
			err := vm.push(err)
			if err != nil {
				return err
			}
		}
	case compiler.OpOpenClass:
		outerCtx := vm.ctx
		newTarget := vm.pop()

		vm.ctx = &object.Context{
			Outer:            outerCtx,
			ExecutionTarget:  newTarget,
			DefinitionTarget: newTarget,
		}
	case compiler.OpCloseClass:
		to := vm.ctx.Outer

		vm.ctx = to
	case compiler.OpSetExecutionTarget:
		oldContext := vm.ctx
		newTarget := vm.pop()

		vm.ctx = &object.Context{
			Outer:            oldContext,
			ExecutionTarget:  newTarget,
			DefinitionTarget: oldContext.DefinitionTarget,
		}
	case compiler.OpResetExecutionTarget:
		vm.ctx = vm.ctx.Outer
	case compiler.OpCloseBlock:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		numFreeVars := compiler.ReadUint8(ins[ip+3:])
		vm.currentFrame().ip += 3

		err := vm.closeBlock(int(constIndex), int(numFreeVars))
		if err != nil {
			return err
		}
	case compiler.OpDefinitionStaticTrue:
		vm.ctx.DefinitionTarget = vm.ctx.DefinitionTarget.Class()
	case compiler.OpDefinitionStaticFalse:
		vm.ctx.DefinitionTarget = vm.ctx.DefinitionTarget.(*object.SingletonClass).Instance
	default:
		if opString, ok := infixOperators[op]; ok {
			left := vm.pop()

			result, sendErr := left.SEND(vm.ctx, vm.evalBlock, opString, left, nil, vm.StackTop())
			if sendErr != nil {
				vm.stack[vm.sp-1] = core.NewStandardError(sendErr.Error())
			} else {
				vm.stack[vm.sp-1] = result
			}
		}
	}

	return nil
}

func (vm *VM) closeBlock(constIndex, numFreeVars int) error {
	constant := vm.constants[constIndex]
	block, ok := constant.(*object.Block)
	if !ok {
		return fmt.Errorf("not a block: %+v", constant)
	}

	free := make([]object.EmeraldValue, numFreeVars)
	for i := 0; i < numFreeVars; i++ {
		free[i] = vm.stack[vm.sp-numFreeVars+i]
	}

	vm.sp = vm.sp - numFreeVars

	return vm.push(object.NewClosedBlock(block, free))
}

func (vm *VM) callFunction(numArgs int) (err object.EmeraldValue) {
	name := vm.stack[vm.sp-2-numArgs].(*core.SymbolInstance)
	block := vm.stack[vm.sp-1-numArgs]

	target := vm.ctx.ExecutionTarget
	method, errVal := target.Class().ExtractMethod(name.Value, target.Class(), target)
	if errVal != nil {
		return core.NewStandardError(errVal.Error())
	}

	switch method := method.(type) {
	case *object.ClosedBlock:
		if numArgs != method.NumArgs {
			return core.NewArgumentError(numArgs, method.NumArgs)
		}

		frame := NewFrame(method, vm.sp-numArgs)
		vm.pushFrame(frame)
		vm.sp = frame.basePointer + method.NumLocals
	case *object.WrappedBuiltInMethod:
		result := method.Method(vm.ctx, target, block, vm.evalBlock, vm.stack[vm.sp-numArgs:vm.sp]...)
		vm.sp -= numArgs + 2
		err := vm.push(result)

		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (vm *VM) Context() *object.Context {
	return vm.ctx
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
	o := vm.StackTop()
	vm.sp--
	return o
}

func (vm *VM) buildArray(startIndex, endIndex int) object.EmeraldValue {
	elements := make([]object.EmeraldValue, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return core.NewArray(elements)
}

func (vm *VM) conditionalJump(condition bool, ins compiler.Instructions, ip int) {
	vm.currentFrame().ip += 2

	if condition {
		newPosition := int(compiler.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip = newPosition - 1
		vm.sp--
	}
}

func isTruthy(obj object.EmeraldValue) bool {
	switch obj {
	case core.FALSE, core.NULL:
		return false
	default:
		return true
	}
}
