package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/kernel"
	"emerald/object"
	"fmt"
)

const (
	GlobalsSize = 65536
)

// VM is our virtual machine responsible for the fetch, decode, execute cycle
type VM struct {
	ctx         *object.Context
	stack       []object.EmeraldValue
	sp          int // Always points to the next value. Top of stack is stack[sp-1]
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
		stack:       make([]object.EmeraldValue, StackSize),
		sp:          0,
		frames:      frames,
		framesIndex: 1,
	}

	for _, option := range options {
		option(vm)
	}

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
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	switch op {
	case compiler.OpPop:
		vm.pop()
	case compiler.OpTrue:
		vm.push(core.TRUE)
	case compiler.OpFalse:
		vm.push(core.FALSE)
	case compiler.OpNull:
		vm.push(core.NULL)
	case compiler.OpPushConstant:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		vm.push(kernel.GetConst(constIndex))
	case compiler.OpJump:
		pos := int(compiler.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip = pos - 1
	case compiler.OpJumpNotTruthy:
		vm.conditionalJump(!core.IsTruthy(vm.StackTop()), ins, ip)
	case compiler.OpJumpTruthy:
		vm.conditionalJump(core.IsTruthy(vm.StackTop()), ins, ip)
	case compiler.OpGetGlobal:
		globalIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2
		vm.push(kernel.GetGlobalVariable(globalIndex))
	case compiler.OpSetGlobal:
		globalIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2
		kernel.SetGlobalVariable(globalIndex, vm.StackTop())
	case compiler.OpGetLocal:
		localIndex := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1
		frame := vm.currentFrame()
		vm.push(vm.stack[frame.basePointer+int(localIndex)])
	case compiler.OpSetLocal:
		localIndex := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1
		frame := vm.currentFrame()
		vm.stack[frame.basePointer+int(localIndex)] = vm.StackTop()
	case compiler.OpGetFree:
		freeIndex := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1

		vm.push(vm.currentFrame().block.FreeVariables[freeIndex])
	case compiler.OpInstanceVarGet:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		name := kernel.GetConst(constIndex)
		target := vm.ctx.ExecutionTarget

		val := target.InstanceVariableGet(name.(*core.SymbolInstance).Value, target, target)

		if val == nil {
			val = core.NULL
		}

		vm.push(val)
	case compiler.OpInstanceVarSet:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		name := kernel.GetConst(constIndex)
		val := vm.StackTop()
		target := vm.ctx.ExecutionTarget

		target.InstanceVariableSet(name.(*core.SymbolInstance).Value, val)
	case compiler.OpArray:
		numElements := int(compiler.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip += 2

		array := vm.buildArray(vm.sp-numElements, vm.sp)
		vm.sp = vm.sp - numElements

		vm.push(array)
	case compiler.OpHash:
		numElements := int(compiler.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip += 2

		startIndex := vm.sp - numElements

		hash := vm.buildHash(startIndex, vm.sp)
		vm.sp = startIndex

		vm.push(hash)
	case compiler.OpBang:
		vm.executeBangOperator()
	case compiler.OpMinus:
		vm.executeMinusOperator()
	case compiler.OpReturn:
		frame := vm.popFrame()
		vm.sp = frame.basePointer - 2

		vm.push(core.NULL)
	case compiler.OpReturnValue:
		returnValue := vm.pop()

		frame := vm.popFrame()
		vm.sp = frame.basePointer - 2

		vm.push(returnValue)
	case compiler.OpDefineMethod:
		block := vm.pop().(*object.Block)
		name := vm.stack[vm.sp-1].(*core.SymbolInstance)

		vm.ctx.DefinitionTarget.DefineMethod(object.NewClosedBlock(block, []object.EmeraldValue{}), name)
	case compiler.OpSend:
		numArgs := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1
		vm.callFunction(int(numArgs))
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

		vm.closeBlock(int(constIndex), int(numFreeVars))
	case compiler.OpDefinitionStaticTrue:
		vm.ctx.DefinitionTarget = vm.ctx.DefinitionTarget.Class()
	case compiler.OpDefinitionStaticFalse:
		vm.ctx.DefinitionTarget = vm.ctx.DefinitionTarget.(*object.SingletonClass).Instance
	default:
		if opString, ok := infixOperators[op]; ok {
			left := vm.pop()

			result, sendErr := left.SEND(vm.ctx, vm.EvalBlock, opString, left, nil, vm.StackTop())
			if sendErr != nil {
				vm.stack[vm.sp-1] = core.NewStandardError(sendErr.Error())
			} else {
				vm.stack[vm.sp-1] = result
			}
		} else {
			def, err := compiler.Lookup(byte(op))
			if err != nil {
				return err
			}

			return fmt.Errorf("opcode not implemented %s", def.Name)
		}
	}

	return err
}

func (vm *VM) closeBlock(constIndex, numFreeVars int) {
	constant := kernel.GetConst(uint16(constIndex))
	block, ok := constant.(*object.Block)
	if !ok {
		panic(fmt.Errorf("not a block: %+v", constant))
	}

	free := make([]object.EmeraldValue, numFreeVars)
	for i := 0; i < numFreeVars; i++ {
		free[i] = vm.stack[vm.sp-numFreeVars+i]
	}

	vm.sp = vm.sp - numFreeVars

	vm.push(object.NewClosedBlock(block, free))
}

func (vm *VM) callFunction(numArgs int) {
	basePointer := vm.sp - numArgs

	name := vm.stack[basePointer-2].(*core.SymbolInstance)
	block := vm.stack[basePointer-1]

	target := vm.ctx.ExecutionTarget
	method, err := target.Class().ExtractMethod(name.Value, target.Class(), target)
	if err != nil {
		panic(err)
	}

	switch method := method.(type) {
	case *object.ClosedBlock:
		frame := NewFrame(method, basePointer)
		vm.pushFrame(frame)
		vm.sp = frame.basePointer + method.NumLocals
	case *object.WrappedBuiltInMethod:
		result := vm.evalBuiltIn(method, block, vm.stack[vm.sp-numArgs:vm.sp])
		vm.sp -= numArgs + 2
		vm.push(result)
	}
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
func (vm *VM) push(obj object.EmeraldValue) {
	if vm.sp >= StackSize {
		panic(fmt.Errorf("stack overflow: max stack size of %d exceeded", StackSize))
	}

	vm.stack[vm.sp] = obj
	vm.sp++
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

func (vm *VM) buildHash(startIndex, endIndex int) object.EmeraldValue {
	hash := core.NewHash()

	for i := startIndex; i < endIndex; i += 2 {
		hash.Set(vm.stack[i], vm.stack[i+1])
	}

	return hash
}

func (vm *VM) conditionalJump(condition bool, ins compiler.Instructions, ip int) {
	vm.currentFrame().ip += 2

	if condition {
		newPosition := int(compiler.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip = newPosition - 1
		vm.sp--
	}
}
