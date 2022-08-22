package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"log"
)

func (vm *VM) evalBuiltIn(builtin *object.WrappedBuiltInMethod, block object.EmeraldValue, args []object.EmeraldValue) object.EmeraldValue {
	return builtin.Method(vm.ctx, vm.ctx.ExecutionTarget, block, vm.EvalBlock, args...)
}

func (vm *VM) EvalBlock(block object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	return vm.withExecutionContextForBlock(func() object.EmeraldValue {
		switch bl := block.(type) {
		case *object.WrappedBuiltInMethod:
			// Builtin methods are easy, just call some Go code
			return vm.evalBuiltIn(bl, core.NULL, args)
		case *object.ClosedBlock:
			var err error

			// The VM accounts for the name of the method being called being on the stack when a block is evaluated
			// So we just push something on the stack and nil is the cheapest
			vm.push(core.NULL)
			// Same for a block value
			vm.push(core.NULL)

			// Add the arguments to the stack
			for _, arg := range args {
				vm.push(arg)
			}

			// Prepare the call frame
			basePointer := vm.sp - len(args)
			startFrameIndex := vm.framesIndex
			vm.pushFrame(NewFrame(bl, basePointer))

			// Prepare the vm stack pointer
			vm.sp = basePointer + bl.NumLocals

			var (
				ip  int
				ins compiler.Instructions
				op  compiler.Opcode
			)

			// Execute
			for vm.framesIndex > startFrameIndex {
				vm.currentFrame().ip++

				ip = vm.currentFrame().ip
				ins = vm.currentFrame().Instructions()
				op = compiler.Opcode(ins[ip])

				err = vm.execute(ip, ins, op)
				if err != nil {
					return core.NewStandardError(err.Error())
				}
			}

			// Return value is left on the stack
			return vm.pop()
		default:
			log.Panicf("Yielded to not a block?, got=%#v", bl)
		}

		return core.NULL
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
