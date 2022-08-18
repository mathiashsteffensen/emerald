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
			return vm.evalBuiltIn(bl, core.NULL, args)
		case *object.ClosedBlock:
			var err error

			err = vm.push(core.NULL)
			if err != nil {
				return core.NewStandardError(err.Error())
			}
			err = vm.push(core.NULL)
			if err != nil {
				return core.NewStandardError(err.Error())
			}

			for _, arg := range args {
				err = vm.push(arg)
				if err != nil {
					return core.NewStandardError(err.Error())
				}
			}

			basePointer := vm.sp - len(args)
			startFrameIndex := vm.framesIndex
			vm.pushFrame(NewFrame(bl, basePointer))
			vm.sp = basePointer + bl.NumLocals

			var (
				ip  int
				ins compiler.Instructions
				op  compiler.Opcode
			)

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
