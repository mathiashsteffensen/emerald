package vm

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/debug"
	"emerald/heap"
	"emerald/object"
	"fmt"
	"os"
)

// VM is our virtual machine responsible for the fetch, decode, execute cycle
type VM struct {
	ctx        *object.Context
	fibers     []*Fiber
	fiberIndex int
}

func New(file string, bytecode *bytecode.Bytecode) *VM {
	mainBlock := &object.ClosedBlock{Block: &object.Block{Bytecode: *bytecode}}
	mainFrame := NewFrame(mainBlock, 0)

	rootFiber := NewFiber(mainFrame)

	vm := &VM{
		fibers:     []*Fiber{rootFiber},
		fiberIndex: 0,
	}

	vm.ctx = vm.newContext(file, core.MainObject, core.NULL)

	heap.SetGlobalVariableString("$LOAD_PATH", core.NewArray([]object.EmeraldValue{
		core.NewString(debug.BinaryDir),
	}))

	argv := make([]object.EmeraldValue, len(os.Args))
	for i, arg := range os.Args {
		argv[i] = core.NewString(arg)
	}
	heap.SetGlobalVariableString("$*", core.NewArray(argv))
	setConst(core.MainObject, "ARGV", core.NewArray(argv))

	object.EvalBlock = func(block *object.ClosedBlock, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return vm.withExecutionContextForBlock(block, func() object.EmeraldValue {
			return vm.rawEvalBlock(block, core.NULL, kwargs, args...)
		})
	}
	core.Send = vm.Send

	return vm
}

func (vm *VM) Run() {
	vm.runWhile(func() bool {
		poppedLastFrame := vm.currentFiber().framesIndex == 0
		frameHasMoreInstructions := func() bool {
			return vm.currentFiber().currentFrame().ip < len(vm.currentFiber().currentFrame().Instructions())-1
		}

		return !poppedLastFrame && frameHasMoreInstructions()
	})
}

func (vm *VM) runWhile(condition func() bool) {
	var (
		ip  int
		ins bytecode.Instructions
		op  bytecode.Opcode
	)

	for condition() {
		if !vm.currentFiber().inRescue && vm.ExceptionIsRaised() {
			rescued := vm.popFramesUntilExceptionRescuedOrProgramTerminates()
			if !rescued {
				break
			}
		}

		vm.currentFiber().currentFrame().ip++

		ip, ins, op = vm.fetch()

		vm.execute(ip, ins, op)
	}
}

func (vm *VM) fetch() (int, bytecode.Instructions, bytecode.Opcode) {
	ip := vm.currentFiber().currentFrame().ip
	ins := vm.currentFiber().currentFrame().Instructions()
	return ip, ins, bytecode.Opcode(ins[ip])
}

func (vm *VM) execute(ip int, ins bytecode.Instructions, op bytecode.Opcode) {
	switch op {
	case bytecode.OpPop:
		vm.pop()
	case bytecode.OpSelf:
		vm.push(vm.ctx.Self)
	case bytecode.OpTrue:
		vm.push(core.TRUE)
	case bytecode.OpFalse:
		vm.push(core.FALSE)
	case bytecode.OpNull:
		vm.push(core.NULL)
	case bytecode.OpStringJoin:
		vm.executeOpStringJoin(ins, ip)
	case bytecode.OpYield:
		vm.executeOpYield(ins, ip)
	case bytecode.OpPushConstant:
		constIndex := vm.readUint16(ins, ip)

		vm.push(heap.GetConstant(constIndex))
	case bytecode.OpAdd:
		vm.evalInfixOperator("+")
	case bytecode.OpSub:
		vm.evalInfixOperator("-")
	case bytecode.OpDiv:
		vm.evalInfixOperator("/")
	case bytecode.OpMul:
		vm.evalInfixOperator("*")
	case bytecode.OpMatch:
		vm.evalInfixOperator("=~")
	case bytecode.OpSpaceship:
		vm.evalInfixOperator("<=>")
	case bytecode.OpLessThan:
		vm.evalInfixOperator("<")
	case bytecode.OpLessThanOrEq:
		vm.evalInfixOperator("<=")
	case bytecode.OpGreaterThan:
		vm.evalInfixOperator(">")
	case bytecode.OpGreaterThanOrEq:
		vm.evalInfixOperator(">=")
	case bytecode.OpEqual:
		vm.evalInfixOperator("==")
	case bytecode.OpCaseEqual:
		vm.evalInfixOperator("===")
	case bytecode.OpNotEqual:
		vm.evalInfixOperator("!=")
	case bytecode.OpBinShiftLeft:
		vm.evalInfixOperator("<<")
	case bytecode.OpJump:
		pos := int(bytecode.ReadUint16(ins[ip+1:]))
		vm.currentFiber().currentFrame().ip = pos - 1
	case bytecode.OpJumpNotTruthy:
		vm.conditionalJump(!core.IsTruthy(vm.StackTop()), ins, ip)
	case bytecode.OpJumpTruthy:
		vm.conditionalJump(core.IsTruthy(vm.StackTop()), ins, ip)
	case bytecode.OpCheckCaseEqual:
		vm.executeOpCheckCaseEqual(ins, ip)
	case bytecode.OpGetGlobal:
		globalIndex := vm.readUint16(ins, ip)
		value := heap.GetGlobalVariable(globalIndex)
		if value == nil {
			value = core.NULL
		}
		vm.push(value)
	case bytecode.OpSetGlobal:
		globalIndex := vm.readUint16(ins, ip)
		heap.SetGlobalVariable(globalIndex, vm.StackTop())
	case bytecode.OpGetLocal:
		localIndex := vm.readUint8(ins, ip)
		frame := vm.currentFiber().currentFrame()
		vm.push(vm.stack()[frame.basePointer+int(localIndex)])
	case bytecode.OpSetLocal:
		localIndex := bytecode.ReadUint8(ins[ip+1:])
		vm.currentFiber().currentFrame().ip += 1
		frame := vm.currentFiber().currentFrame()
		vm.stack()[frame.basePointer+int(localIndex)] = vm.StackTop()
	case bytecode.OpGetFree:
		freeIndex := bytecode.ReadUint8(ins[ip+1:])
		vm.currentFiber().currentFrame().ip += 1

		vm.push(vm.currentFiber().currentFrame().block.FreeVariables[freeIndex])
	case bytecode.OpInstanceVarGet:
		constIndex := vm.readUint16(ins, ip)

		name := heap.GetConstant(constIndex)
		target := vm.ctx.Self

		val := target.InstanceVariableGet(name.(*core.SymbolInstance).Value, target, target)

		if val == nil {
			val = core.NULL
		}

		vm.push(val)
	case bytecode.OpConstantGet:
		vm.executeOpConstantGet(ins, ip)
	case bytecode.OpConstantSet:
		vm.executeOpConstantSet(ins, ip)
	case bytecode.OpScopedConstantGet:
		vm.executeOpScopedConstantGet(ins, ip)
	case bytecode.OpInstanceVarSet:
		constIndex := bytecode.ReadUint16(ins[ip+1:])
		vm.currentFiber().currentFrame().ip += 2

		name := heap.GetConstant(constIndex)
		val := vm.StackTop()
		target := vm.ctx.Self

		target.InstanceVariableSet(name.(*core.SymbolInstance).Value, val)
	case bytecode.OpArray:
		numElements := int(bytecode.ReadUint16(ins[ip+1:]))
		vm.currentFiber().currentFrame().ip += 2

		array := vm.buildArray(vm.currentFiber().sp-numElements, vm.currentFiber().sp)
		vm.currentFiber().sp = vm.currentFiber().sp - numElements

		vm.push(array)
	case bytecode.OpHash:
		numElements := int(bytecode.ReadUint16(ins[ip+1:]))
		vm.currentFiber().currentFrame().ip += 2

		startIndex := vm.currentFiber().sp - numElements

		hash := vm.buildHash(startIndex, vm.currentFiber().sp)
		vm.currentFiber().sp = startIndex

		vm.push(hash)
	case bytecode.OpBang:
		vm.executeBangOperator()
	case bytecode.OpMinus:
		vm.executeOpMinus()
	case bytecode.OpReturn:
		vm.executeOpReturn()
	case bytecode.OpReturnValue:
		vm.executeOpReturnValue()
	case bytecode.OpDefineMethod:
		block := vm.pop().(*object.Block)
		name := vm.stack()[vm.currentFiber().sp-1].(*core.SymbolInstance)

		vm.ctx.Self.DefinedMethodSet()[name.Value] = object.NewClosedBlock(nil, block, []object.EmeraldValue{}, vm.ctx.File, vm.ctx.DefaultMethodVisibility)
	case bytecode.OpSend:
		numArgs := vm.readUint8(ins, ip)
		hasKwargs := vm.readUint8(ins, ip+1)
		vm.callMethod(int(numArgs), hasKwargs == 1)
	case bytecode.OpOpenClass:
		// Fetch the symbol name from the heap
		nameIndex := vm.readUint16(ins, ip)
		name := heap.GetConstant(nameIndex).(*core.SymbolInstance).Value

		// Parent class is top off stack
		parent := vm.pop()

		// Check if the class is already defined
		class, err := getConst(vm.ctx.Self, name)
		if err != nil {
			// If not create a new class object
			class = core.DefineNestedClass(vm.ctx.Self, name, parent.(*object.Class))
		}

		// Create and set a new Context with this class as Self
		vm.ctx = vm.newEnclosedContext(vm.ctx.File, class, vm.ctx.Block)
	case bytecode.OpOpenModule:
		outerCtx := vm.ctx
		nameIndex := vm.readUint16(ins, ip)
		name := heap.GetConstant(nameIndex).(*core.SymbolInstance).Value

		module, err := getConst(vm.ctx.Self, name)
		if err != nil {
			module = core.DefineNestedModule(vm.ctx.Self, name)
		}

		vm.ctx = vm.newEnclosedContext(outerCtx.File, module, outerCtx.Block)
	case bytecode.OpUnwrapContext:
		to := vm.ctx.Outer

		vm.ctx = to
	case bytecode.OpCloseBlock:
		constIndex := bytecode.ReadUint16(ins[ip+1:])
		numFreeVars := bytecode.ReadUint8(ins[ip+3:])
		vm.currentFiber().currentFrame().ip += 3

		vm.closeBlock(int(constIndex), int(numFreeVars))
	case bytecode.OpStaticTrue:
		vm.ctx.Self = vm.ctx.Self.Class()
	case bytecode.OpStaticFalse:
		vm.ctx.Self = vm.ctx.Self.(*object.SingletonClass).Instance
	default:
		def, err := bytecode.Lookup(byte(op))
		if err != nil {
			panic(err)
		}

		panic(fmt.Errorf("opcode not implemented %s", def.Name))
	}
}

func (vm *VM) closeBlock(constIndex, numFreeVars int) {
	constant := heap.GetConstant(uint16(constIndex))
	block := constant.(*object.Block)

	free := make([]object.EmeraldValue, numFreeVars)
	for i := 0; i < numFreeVars; i++ {
		free[i] = vm.stack()[vm.currentFiber().sp-numFreeVars+i]
	}

	vm.currentFiber().sp = vm.currentFiber().sp - numFreeVars

	vm.push(object.NewClosedBlock(vm.ctx, block, free, "", object.PUBLIC))
}

func (vm *VM) callMethod(numArgs int, hasKwargs bool) {
	var (
		kwargsHash   *core.HashInstance
		kwargsMap    = map[string]object.EmeraldValue{}
		basePointer  int
		argsEndIndex int
	)

	if hasKwargs {
		var ok bool
		kwargsHash, ok = vm.currentFiber().pop().(*core.HashInstance)
		if !ok {
			debug.FatalBugF("Keyword arguments instance was not a hash? got %s", vm.currentFiber().StackTop().Inspect())
		}
		basePointer = vm.currentFiber().sp - (numArgs - kwargsHash.Values.Len())
	} else {
		basePointer = vm.currentFiber().sp - numArgs
	}

	receiver := vm.stack()[basePointer-3]
	name, ok := vm.stack()[basePointer-2].(*core.SymbolInstance)
	if !ok {
		debug.FatalBugF("Method name instance was not a symbol? got %q", vm.stack()[basePointer-2].Inspect())
	}
	block := vm.stack()[basePointer-1]

	method, visibility, isDefinedOnReceiver, err := receiver.Class().ExtractMethod(name.Value, receiver.Class(), receiver)
	if err != nil {
		raiseUndefinedNoMethodError(name.Value, receiver)
	}

	if ok := vm.ctx.ValidateMethodVisibility(receiver, visibility, isDefinedOnReceiver); !ok {
		raiseNotVisibleNoMethodError(name.Value, receiver)
	}

	// Handy for debugging, but makes the VM quite slow when calling DebugF in a hot path
	// debug.DebugF("Calling method %s#%s %d %s", receiver.Inspect(), name.Value, numArgs)

	vm.withExecutionContext(receiver, block, func() {
		switch method := method.(type) {
		case *object.ClosedBlock:
			if hasKwargs {
				sortedKwargsHash := core.NewHash()

				// Sort kwargs first, so they match the definition order, this allows local variable references to resolve correctly
				for _, kwargStringKey := range method.Kwargs {
					symbolKey := core.NewSymbol(kwargStringKey)

					sortedKwargsHash.Set(symbolKey, kwargsHash.Get(symbolKey))
				}

				kwargsMap, argsEndIndex = vm.pushKwargsToStack(sortedKwargsHash)
			} else {
				argsEndIndex = vm.currentFiber().sp
			}

			if _, err := core.EnforceArity(vm.stack()[basePointer:argsEndIndex], kwargsMap, method.NumArgs, method.NumArgs, method.Kwargs...); err != nil {
				return
			}

			frame := NewFrame(method, basePointer)
			vm.currentFiber().pushFrame(frame)
			vm.currentFiber().sp = frame.basePointer + method.NumLocals
			originalFrameIndex := vm.currentFiber().framesIndex

			if method.File != "" {
				vm.ctx.File = method.File
			}

			vm.runWhile(func() bool {
				return vm.currentFiber().framesIndex >= originalFrameIndex
			})
		case *object.WrappedBuiltInMethod:
			if hasKwargs {
				kwargsMap, argsEndIndex = vm.pushKwargsToStack(kwargsHash)
			} else {
				argsEndIndex = vm.currentFiber().sp
			}

			result := vm.evalBuiltIn(method, block, vm.stack()[basePointer:argsEndIndex], kwargsMap)
			vm.currentFiber().sp = basePointer - 3
			vm.push(result)
		}
	})
}

func raiseUndefinedNoMethodError(name string, receiver object.EmeraldValue) {
	core.Raise(
		core.NewNoMethodError(
			fmt.Sprintf("undefined method '%s' for %s:%s", name, receiver.Inspect(), receiver.Class().Super().(*object.Class).Name),
		),
	)
}

func raiseNotVisibleNoMethodError(name string, receiver object.EmeraldValue) {
	var receiverPart string
	receiverClassName := receiver.Class().Super().(*object.Class).Name
	if receiverClassName == core.Class.Name {
		receiverPart = fmt.Sprintf("%s:%s", receiver.Inspect(), receiverClassName)
	} else {
		receiverPart = receiverClassName
	}

	core.Raise(
		core.NewNoMethodError(
			fmt.Sprintf("private method `%s' called for %s", name, receiverPart),
		),
	)
}

func (vm *VM) pushKwargsToStack(kwargsHash *core.HashInstance) (map[string]object.EmeraldValue, int) {
	kwargsMap := map[string]object.EmeraldValue{}

	kwargsHash.Each(func(key object.EmeraldValue, value object.EmeraldValue) {
		vm.push(value)
		kwargsMap[key.Inspect()] = value
	})

	argsEndIndex := vm.currentFiber().sp - kwargsHash.Values.Len()

	return kwargsMap, argsEndIndex
}

func (vm *VM) evalInfixOperator(op string) {
	left := vm.pop()

	vm.stack()[vm.currentFiber().sp-1] = vm.Send(left, op, core.NULL, map[string]object.EmeraldValue{}, vm.StackTop())
}

func (vm *VM) Context() *object.Context {
	return vm.ctx
}

func (vm *VM) withExecutionContext(self object.EmeraldValue, block object.EmeraldValue, cb func()) {
	oldCtx := vm.ctx

	vm.ctx = vm.newEnclosedContext(oldCtx.File, self, block)

	cb()

	oldCtx.DefaultMethodVisibility = vm.ctx.DefaultMethodVisibility

	vm.ctx = oldCtx
}

// StackTop fetches the object at the top of the stack
func (vm *VM) StackTop() object.EmeraldValue {
	return vm.currentFiber().StackTop()
}

func (vm *VM) LastPoppedStackElem() object.EmeraldValue {
	return vm.stack()[vm.currentFiber().sp]
}

// push an obj on to the stack
func (vm *VM) push(obj object.EmeraldValue) {
	vm.currentFiber().push(obj)
}

// pop an obj from the top of the stack
func (vm *VM) pop() object.EmeraldValue {
	return vm.currentFiber().pop()
}

func (vm *VM) buildArray(startIndex, endIndex int) object.EmeraldValue {
	elements := make([]object.EmeraldValue, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack()[i]
	}

	return core.NewArray(elements)
}

func (vm *VM) buildHash(startIndex, endIndex int) object.EmeraldValue {
	hash := core.NewHash()

	for i := startIndex; i < endIndex; i += 2 {
		hash.Set(vm.stack()[i], vm.stack()[i+1])
	}

	return hash
}

func (vm *VM) conditionalJump(condition bool, ins bytecode.Instructions, ip int) {
	vm.currentFiber().currentFrame().ip += 2

	if condition {
		newPosition := int(bytecode.ReadUint16(ins[ip+1:]))
		vm.currentFiber().currentFrame().ip = newPosition - 1
		vm.currentFiber().sp--
	}
}

func (vm *VM) readUint8(ins bytecode.Instructions, ip int) uint8 {
	val := bytecode.ReadUint8(ins[ip+1:])
	vm.currentFiber().currentFrame().ip += 1
	return val
}

func (vm *VM) readUint16(ins bytecode.Instructions, ip int) uint16 {
	val := bytecode.ReadUint16(ins[ip+1:])
	vm.currentFiber().currentFrame().ip += 2
	return val
}
