package vm

import (
	"emerald/core"
	"emerald/object"
	"log"
)

func (vm *VM) evalBuiltIn(receiver object.EmeraldValue, builtin *object.WrappedBuiltInMethod, block object.EmeraldValue, args []object.EmeraldValue) object.EmeraldValue {
	return builtin.Method(vm.ctx, receiver, block, vm.Yield, args...)
}

func (vm *VM) Yield(block object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	return vm.withExecutionContextForBlock(func() object.EmeraldValue {
		return vm.rawEvalBlock(block, core.NULL, args...)
	})
}

func (vm *VM) withExecutionContextForBlock(cb func() object.EmeraldValue) object.EmeraldValue {
	oldCtx := vm.ctx

	if vm.ctx.Outer != nil {
		vm.ctx = vm.ctx.Outer
	}

	val := cb()

	vm.ctx = oldCtx

	return val
}

func (vm *VM) Send(self object.EmeraldValue, name string, block object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	oldCtx := vm.ctx
	vm.ctx = &object.Context{
		Outer: oldCtx,
		Self:  self,
	}

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
		return vm.evalBuiltIn(vm.ctx.Self, bl, core.NULL, args)
	case *object.ClosedBlock:
		// Method receiver
		vm.push(vm.ctx.Self)
		// The VM accounts for the name of the method being called being on the stack when a method is evaluated
		// So we just push something on the stack and nil is the cheapest
		vm.push(core.NULL)
		// Same for a block value
		vm.push(block)

		// Add the arguments to the stack
		for _, arg := range args {
			vm.push(arg)
		}

		// Prepare the call frame
		startFrameIndex := vm.framesIndex
		basePointer := vm.sp - len(args)
		vm.pushFrame(NewFrame(bl, basePointer))

		// Prepare the vm stack pointer
		vm.sp = basePointer + bl.NumLocals

		// Execute
		vm.runWhile(func() bool {
			return vm.framesIndex > startFrameIndex
		})

		// Return value is left on the stack
		return vm.pop()
	default:
		log.Panicf("Yielded to not a method?, got=%#v", bl)
	}

	return core.NULL
}
