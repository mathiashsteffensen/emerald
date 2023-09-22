package vm

import (
	"emerald/core"
	"emerald/object"
	"fmt"
)

func (vm *VM) Yield(kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	blockToEvaluate := vm.ctx.Block
	return vm.withExecutionContextForBlock(blockToEvaluate, func() object.EmeraldValue {
		return vm.rawEvalBlock(blockToEvaluate, core.NULL, kwargs, args...)
	})
}

func (vm *VM) withExecutionContextForBlock(block object.EmeraldValue, cb func() object.EmeraldValue) object.EmeraldValue {
	oldCtx := vm.ctx

	if closedBlock, ok := block.(*object.ClosedBlock); ok && closedBlock.Context != nil {
		vm.ctx = closedBlock.Context
	}

	val := cb()

	vm.ctx = oldCtx

	return val
}

func (vm *VM) Send(self object.EmeraldValue, name string, block object.EmeraldValue, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	oldCtx := vm.ctx
	vm.ctx = vm.newEnclosedContext(oldCtx.File, self, block)

	method, _, _, err := self.ExtractMethod(name, self.Class(), self)
	if err != nil {
		panic(err)
	}

	result := vm.rawEvalBlock(method, block, kwargs, args...)

	vm.ctx = oldCtx

	return result
}

func (vm *VM) rawEvalBlock(method object.EmeraldValue, block object.EmeraldValue, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	switch bl := method.(type) {
	case *object.WrappedBuiltInMethod:
		// Builtin methods are easy, just call some Go code
		return vm.evalBuiltIn(bl, block, args, kwargs)
	case *object.ClosedBlock:
		if bl.EnforceArity {
			if _, err := core.EnforceArity(args, kwargs, bl.NumArgs, bl.NumArgs, bl.Kwargs...); err != nil {
				return err
			}
		}

		// Method receiver
		vm.push(vm.ctx.Self)

		// The VM accounts for the name of the method being called being on the stack when a method is evaluated
		// So we just push something on the stack
		vm.push(core.NULL)

		vm.push(block)

		// Add the arguments to the stack
		for _, arg := range args {
			vm.push(arg)
		}

		if len(kwargs) != 0 {
			sortedKwargsHash := core.NewHash()

			// Sort kwargs first, so they match the definition order, this allows local variable references to resolve correctly
			for kwargStringKey, value := range kwargs {
				symbolKey := core.NewSymbol(kwargStringKey)

				sortedKwargsHash.Set(symbolKey, value)
			}

			vm.pushKwargsToStack(sortedKwargsHash)
		}

		// Prepare the call frame
		startFrameIndex := vm.currentFiber().framesIndex
		basePointer := vm.currentFiber().sp - len(args) - len(kwargs)
		vm.currentFiber().pushFrame(NewFrame(bl, basePointer))

		// Prepare the vm stack pointer
		vm.currentFiber().sp = basePointer + bl.NumLocals

		// Execute
		vm.runWhile(func() bool {
			return vm.currentFiber().framesIndex > startFrameIndex
		})

		if vm.currentFiber().inRescue || !vm.ExceptionIsRaised() {
			// Return value is left on the stack
			return vm.pop()
		}
	default:
		core.Raise(core.NewException(fmt.Sprintf("yielded to not a method?, got=%s", bl.Inspect())))
	}

	return core.NULL
}

func (vm *VM) evalBuiltIn(builtin *object.WrappedBuiltInMethod, block object.EmeraldValue, args []object.EmeraldValue, kwargs map[string]object.EmeraldValue) object.EmeraldValue {
	oldBlock := vm.ctx.Block
	vm.ctx.Block = block

	result := builtin.Method(vm.ctx, kwargs, args...)

	vm.ctx.Block = oldBlock

	return result
}
