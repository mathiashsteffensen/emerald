package vm

import (
	"emerald/core"
	"emerald/object"
	"fmt"
)

func (vm *VM) Yield(args ...object.EmeraldValue) object.EmeraldValue {
	blockToEvaluate := vm.ctx.Block
	return vm.withExecutionContextForBlock(blockToEvaluate, func() object.EmeraldValue {
		return vm.rawEvalBlock(blockToEvaluate, core.NULL, args...)
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

func (vm *VM) Send(self object.EmeraldValue, name string, block object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	oldCtx := vm.ctx
	vm.ctx = vm.newEnclosedContext(oldCtx.File, self, block)

	method, err := self.ExtractMethod(name, self.Class(), self)
	if err != nil {
		panic(err)
	}

	result := vm.rawEvalBlock(method, block, args...)

	vm.ctx = oldCtx

	return result
}

func (vm *VM) rawEvalBlock(method object.EmeraldValue, block object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	switch bl := method.(type) {
	case *object.WrappedBuiltInMethod:
		// Builtin methods are easy, just call some Go code
		return vm.evalBuiltIn(bl, block, args)
	case *object.ClosedBlock:
		// Method receiver
		vm.push(vm.ctx.Self)

		// The VM accounts for the name of the method being called being on the stack when a method is evaluated
		// So we just push something on the stack and nil is the cheapest
		vm.push(core.NULL)

		vm.push(block)

		// Add the arguments to the stack
		for _, arg := range args {
			vm.push(arg)
		}

		// Prepare the call frame
		startFrameIndex := vm.currentFiber().framesIndex
		basePointer := vm.currentFiber().sp - len(args)
		vm.currentFiber().pushFrame(NewFrame(bl, basePointer))

		// Prepare the vm stack pointer
		vm.currentFiber().sp = basePointer + bl.NumLocals

		// Execute
		vm.runWhile(func() bool {
			return vm.currentFiber().framesIndex > startFrameIndex
		})

		if vm.inRescue || !vm.ExceptionIsRaised() {
			// Return value is left on the stack
			return vm.pop()
		}
	default:
		panic(fmt.Errorf("yielded to not a method?, got=%s", bl.Inspect()))
	}

	return core.NULL
}

func (vm *VM) evalBuiltIn(builtin *object.WrappedBuiltInMethod, block object.EmeraldValue, args []object.EmeraldValue) object.EmeraldValue {
	oldBlock := vm.ctx.Block
	vm.ctx.Block = block

	result := builtin.Method(vm.ctx, args...)

	vm.ctx.Block = oldBlock

	return result
}
