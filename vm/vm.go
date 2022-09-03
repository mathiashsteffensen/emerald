package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/heap"
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

func New(file string, bytecode *compiler.Bytecode) *VM {
	mainBlock := &object.ClosedBlock{Block: &object.Block{Instructions: bytecode.Instructions}}
	mainFrame := NewFrame(mainBlock, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	vm := &VM{
		stack:       make([]object.EmeraldValue, StackSize),
		sp:          0,
		frames:      frames,
		framesIndex: 1,
	}

	vm.ctx = vm.newContext(file, core.MainObject, core.NULL)

	object.EvalBlock = func(block *object.ClosedBlock, args ...object.EmeraldValue) object.EmeraldValue {
		return vm.withExecutionContextForBlock(block, func() object.EmeraldValue {
			return vm.rawEvalBlock(block, core.NULL, args...)
		})
	}
	core.Send = vm.Send

	return vm
}

func (vm *VM) Run() {
	vm.runWhile(func() bool {
		return vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1
	})
}

func (vm *VM) runWhile(condition func() bool) {
	var (
		ip  int
		ins compiler.Instructions
		op  compiler.Opcode
	)

	for condition() {
		vm.currentFrame().ip++

		ip, ins, op = vm.fetch()

		vm.execute(ip, ins, op)
	}
}

func (vm *VM) fetch() (int, compiler.Instructions, compiler.Opcode) {
	ip := vm.currentFrame().ip
	ins := vm.currentFrame().Instructions()
	return ip, ins, compiler.Opcode(ins[ip])
}

func (vm *VM) execute(ip int, ins compiler.Instructions, op compiler.Opcode) {
	switch op {
	case compiler.OpPop:
		vm.pop()
	case compiler.OpSelf:
		vm.push(vm.ctx.Self)
	case compiler.OpTrue:
		vm.push(core.TRUE)
	case compiler.OpFalse:
		vm.push(core.FALSE)
	case compiler.OpNull:
		vm.push(core.NULL)
	case compiler.OpYield:
		vm.executeOpYield(ins, ip)
	case compiler.OpPushConstant:
		constIndex := vm.readUint16(ins, ip)

		vm.push(heap.GetConstant(constIndex))
	case compiler.OpAdd:
		vm.evalInfixOperator("+")
	case compiler.OpSub:
		vm.evalInfixOperator("-")
	case compiler.OpDiv:
		vm.evalInfixOperator("/")
	case compiler.OpMul:
		vm.evalInfixOperator("*")
	case compiler.OpMatch:
		vm.evalInfixOperator("=~")
	case compiler.OpSpaceship:
		vm.evalInfixOperator("<=>")
	case compiler.OpLessThan:
		vm.evalInfixOperator("<")
	case compiler.OpLessThanOrEq:
		vm.evalInfixOperator("<=")
	case compiler.OpGreaterThan:
		vm.evalInfixOperator(">")
	case compiler.OpGreaterThanOrEq:
		vm.evalInfixOperator(">=")
	case compiler.OpEqual:
		vm.evalInfixOperator("==")
	case compiler.OpNotEqual:
		vm.evalInfixOperator("!=")
	case compiler.OpJump:
		pos := int(compiler.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip = pos - 1
	case compiler.OpJumpNotTruthy:
		vm.conditionalJump(!core.IsTruthy(vm.StackTop()), ins, ip)
	case compiler.OpJumpTruthy:
		vm.conditionalJump(core.IsTruthy(vm.StackTop()), ins, ip)
	case compiler.OpGetGlobal:
		globalIndex := vm.readUint16(ins, ip)
		value := heap.GetGlobalVariable(globalIndex)
		if value == nil {
			value = core.NULL
		}
		vm.push(value)
	case compiler.OpSetGlobal:
		globalIndex := vm.readUint16(ins, ip)
		heap.SetGlobalVariable(globalIndex, vm.StackTop())
	case compiler.OpGetLocal:
		localIndex := vm.readUint8(ins, ip)
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
		constIndex := vm.readUint16(ins, ip)

		name := heap.GetConstant(constIndex)
		target := vm.ctx.Self

		val := target.InstanceVariableGet(name.(*core.SymbolInstance).Value, target, target)

		if val == nil {
			val = core.NULL
		}

		vm.push(val)
	case compiler.OpConstantGet:
		vm.executeOpConstantGet(ins, ip)
	case compiler.OpConstantSet:
		vm.executeOpConstantSet(ins, ip)
	case compiler.OpConstantGetOrSet:
		vm.executeOpConstantGetOrSet(ins, ip)
	case compiler.OpScopedConstantGet:
		vm.executeOpScopedConstantGet(ins, ip)
	case compiler.OpInstanceVarSet:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		name := heap.GetConstant(constIndex)
		val := vm.StackTop()
		target := vm.ctx.Self

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
		vm.sp = frame.basePointer - 3

		vm.push(core.NULL)
	case compiler.OpReturnValue:
		returnValue := vm.stack[vm.sp-1]

		frame := vm.popFrame()
		vm.sp = frame.basePointer - 3

		vm.push(returnValue)
	case compiler.OpDefineMethod:
		block := vm.pop().(*object.Block)
		name := vm.stack[vm.sp-1].(*core.SymbolInstance)

		vm.ctx.Self.DefineMethod(object.NewClosedBlock(nil, block, []object.EmeraldValue{}, vm.ctx.File), name)
	case compiler.OpSend:
		numArgs := compiler.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1
		vm.callFunction(int(numArgs))
	case compiler.OpOpenClass:
		outerCtx := vm.ctx
		newSelf := vm.pop()

		var name string

		switch newSelf := newSelf.(type) {
		case *object.Class:
			name = newSelf.Name
		case *object.Module:
			name = newSelf.Name
		}

		setConst(outerCtx.Self, name, newSelf)

		vm.ctx = vm.newEnclosedContext(outerCtx.File, newSelf, outerCtx.Block)
	case compiler.OpCloseClass:
		to := vm.ctx.Outer

		vm.ctx = to
	case compiler.OpSetExecutionTarget:
		oldContext := vm.ctx
		newSelf := vm.pop()

		vm.ctx = vm.newEnclosedContext(oldContext.File, newSelf, oldContext.Block)
	case compiler.OpResetExecutionTarget:
		vm.ctx = vm.ctx.Outer
	case compiler.OpCloseBlock:
		constIndex := compiler.ReadUint16(ins[ip+1:])
		numFreeVars := compiler.ReadUint8(ins[ip+3:])
		vm.currentFrame().ip += 3

		vm.closeBlock(int(constIndex), int(numFreeVars))
	case compiler.OpStaticTrue:
		vm.ctx.Self = vm.ctx.Self.Class()
	case compiler.OpStaticFalse:
		vm.ctx.Self = vm.ctx.Self.(*object.SingletonClass).Instance
	default:
		def, err := compiler.Lookup(byte(op))
		if err != nil {
			panic(err)
		}

		panic(fmt.Errorf("opcode not implemented %s", def.Name))
	}
}

func (vm *VM) closeBlock(constIndex, numFreeVars int) {
	constant := heap.GetConstant(uint16(constIndex))
	block, ok := constant.(*object.Block)
	if !ok {
		panic(fmt.Errorf("not a block: %+v", constant))
	}

	free := make([]object.EmeraldValue, numFreeVars)
	for i := 0; i < numFreeVars; i++ {
		free[i] = vm.stack[vm.sp-numFreeVars+i]
	}

	vm.sp = vm.sp - numFreeVars

	vm.push(object.NewClosedBlock(vm.ctx, block, free, ""))
}

func (vm *VM) callFunction(numArgs int) {
	basePointer := vm.sp - numArgs

	receiver := vm.stack[basePointer-3]
	name := vm.stack[basePointer-2].(*core.SymbolInstance)
	block := vm.stack[basePointer-1]

	method, err := receiver.Class().ExtractMethod(name.Value, receiver.Class(), receiver)
	if err != nil {
		panic(err)
	}

	vm.withExecutionContext(receiver, block, func() {
		switch method := method.(type) {
		case *object.ClosedBlock:
			frame := NewFrame(method, basePointer)
			vm.pushFrame(frame)
			vm.sp = frame.basePointer + method.NumLocals
			originalFrameIndex := vm.framesIndex

			if method.File != "" {
				vm.ctx.File = method.File
			}

			vm.runWhile(func() bool {
				return vm.framesIndex >= originalFrameIndex
			})
		case *object.WrappedBuiltInMethod:
			result := vm.evalBuiltIn(method, block, vm.stack[basePointer:vm.sp])
			vm.sp = basePointer - 3
			vm.push(result)
		}
	})
}

func (vm *VM) evalInfixOperator(op string) {
	left := vm.pop()

	var result object.EmeraldValue

	vm.withExecutionContext(left, core.NULL, func() {
		result = left.SEND(vm.ctx, op, vm.StackTop())
	})

	vm.stack[vm.sp-1] = result
}

func (vm *VM) Context() *object.Context {
	return vm.ctx
}

func (vm *VM) withExecutionContext(self object.EmeraldValue, block object.EmeraldValue, cb func()) {
	oldCtx := vm.ctx

	vm.ctx = vm.newEnclosedContext(oldCtx.File, self, block)

	cb()

	vm.ctx = oldCtx
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

func (vm *VM) readUint8(ins compiler.Instructions, ip int) uint8 {
	val := compiler.ReadUint8(ins[ip+1:])
	vm.currentFrame().ip += 1
	return val
}

func (vm *VM) readUint16(ins compiler.Instructions, ip int) uint16 {
	val := compiler.ReadUint16(ins[ip+1:])
	vm.currentFrame().ip += 2
	return val
}
