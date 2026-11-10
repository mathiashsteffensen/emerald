package vm

import (
	"emerald/core"
	"emerald/object"
	"fmt"
	"sort"
	"strings"
)

func (vm *VM) Yield(kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	blockToEvaluate := vm.ctx.Block
	return vm.withExecutionContextForBlock(blockToEvaluate, func() object.EmeraldValue {
		return vm.rawEvalBlock(blockToEvaluate, vm.rt.NULL, kwargs, args...)
	})
}

func (vm *VM) withExecutionContextForBlock(block object.EmeraldValue, cb func() object.EmeraldValue) object.EmeraldValue {
	oldCtx := vm.ctx
	defer func() {
		vm.ctx = oldCtx
	}()

	if closedBlock, ok := block.Heap.(*object.ClosedBlock); ok && closedBlock.Context != nil {
		vm.ctx = closedBlock.Context
		// A Proc can outlive the Engine.Eval that created it. Point the
		// captured context at this VM so yield keeps working.
		vm.ctx.ExecutionContext = vm.executionContext
		vm.ctx.Yield = vm.Yield
		vm.ctx.BlockGiven = vm.BlockGiven
	}

	return cb()
}

func (vm *VM) Send(self object.EmeraldValue, name string, block object.EmeraldValue, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	oldCtx := vm.ctx
	vm.ctx = vm.newEnclosedContext(oldCtx.File, self, block)
	defer func() {
		vm.ctx = oldCtx
	}()

	method, err := vm.extractMethod(self, name)
	if err != nil {
		return object.NewHeapObject(err)
	}

	result := vm.rawEvalBlock(method, block, kwargs, args...)

	return result
}

func (vm *VM) rawEvalBlock(method object.EmeraldValue, block object.EmeraldValue, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	switch bl := method.Heap.(type) {
	case *object.WrappedBuiltInMethod:
		// Builtin methods are easy, just call some Go code
		return vm.evalBuiltIn(bl, block, args, kwargs)
	case *object.ClosedBlock:
		if bl.EnforceArity {
			if _, err := vm.rt.EnforceArity(args, kwargs, bl.MinArgs, bl.NumArgs, bl.Kwargs...); err != nil {
				return object.NewHeapObject(err)
			}
		} else if len(bl.Kwargs) != 0 {
			if err := vm.rt.RequireKeywords(kwargs, bl.Kwargs...); err != nil {
				return object.NewHeapObject(err)
			}
		}

		// Method receiver
		vm.push(vm.ctx.Self)

		vm.push(block)

		basePointer := vm.currentFiber().sp

		var argGiven []bool
		if bl.EnforceArity {
			// Add the arguments to the stack
			for _, arg := range args {
				vm.push(arg)
			}

			argGiven = vm.bindPositionalArgs(bl, basePointer, len(args))
			vm.currentFiber().sp = basePointer + bl.NumArgs
		} else {
			vm.bindBlockArguments(bl, basePointer, kwargs, args)
		}

		if bl.EnforceArity && len(kwargs) != 0 {
			sortedKwargsHashVal := vm.rt.NewHash()
			sortedKwargsHash := sortedKwargsHashVal.Heap.(*core.HashInstance)

			// Sort kwargs first, so they match the definition order, this allows local variable references to resolve correctly
			for _, name := range bl.Kwargs {
				sortedKwargsHash.Set(vm.rt.NewSymbol(name), kwargs[":"+name])
			}

			vm.pushKwargsToStack(sortedKwargsHash)
		}

		// Prepare the call frame
		startFrameIndex := vm.currentFiber().framesIndex
		frame := NewFrame(bl, basePointer)
		frame.argGiven = argGiven
		vm.currentFiber().pushFrame(frame)

		// Prepare the vm stack pointer
		vm.currentFiber().sp = basePointer + bl.NumLocals

		// Execute
		vm.runWhile(func() bool {
			return vm.currentFiber().framesIndex > startFrameIndex
		})

		if vm.executionError == nil && !vm.ExceptionIsRaised() {
			// Return value is left on the stack
			return vm.pop()
		}
	default:
		vm.rt.Raise(vm.rt.NewException(fmt.Sprintf("yielded to not a method?, got=%s", method.Inspect())))
	}

	return vm.rt.NULL
}

// bindBlockArguments lays out a block call's arguments in its local slots.
// Blocks are lenient about positional arity: a parameter without an argument
// is nil and extra arguments are dropped. Declared keyword parameters take
// their value by name at the declared slot. With no declared keyword
// parameters, Ruby passes the keywords on as one trailing positional Hash.
// Remaining locals start as nil.
func (vm *VM) bindBlockArguments(bl *object.ClosedBlock, basePointer int, kwargs map[string]object.EmeraldValue, args []object.EmeraldValue) {
	positional := args
	if len(bl.Kwargs) == 0 && len(kwargs) != 0 {
		positional = make([]object.EmeraldValue, len(args), len(args)+1)
		copy(positional, args)
		positional = append(positional, vm.kwargsHash(kwargs))
	}

	for i := 0; i < bl.NumArgs; i++ {
		if i < len(positional) {
			vm.push(positional[i])
		} else {
			vm.push(vm.rt.NULL)
		}
	}

	for _, name := range bl.Kwargs {
		vm.push(kwargs[":"+name])
	}

	for vm.currentFiber().sp < basePointer+bl.NumLocals {
		vm.push(vm.rt.NULL)
	}
}

// kwargsHash builds the Hash Ruby passes to a block that declares no keyword
// parameters. Sorted keys keep the Hash order stable.
func (vm *VM) kwargsHash(kwargs map[string]object.EmeraldValue) object.EmeraldValue {
	hashVal := vm.rt.NewHash()
	hash := hashVal.Heap.(*core.HashInstance)

	keys := make([]string, 0, len(kwargs))
	for key := range kwargs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		hash.Set(vm.rt.NewSymbol(strings.TrimPrefix(key, ":")), kwargs[key])
	}

	return hashVal
}

func (vm *VM) evalBuiltIn(builtin *object.WrappedBuiltInMethod, block object.EmeraldValue, args []object.EmeraldValue, kwargs map[string]object.EmeraldValue) object.EmeraldValue {
	if !vm.checkContext() {
		return vm.rt.NULL
	}

	fiber := vm.currentFiber()
	fiber.nativeFrames = append(fiber.nativeFrames, nativeFrame{name: builtin.Name, depth: fiber.framesIndex})

	oldBlock := vm.ctx.Block
	vm.ctx.Block = block
	defer func() {
		vm.ctx.Block = oldBlock
		fiber.nativeFrames = fiber.nativeFrames[:len(fiber.nativeFrames)-1]
	}()

	result := builtin.Method(vm.ctx, kwargs, args...)
	vm.checkContext()

	return result
}
