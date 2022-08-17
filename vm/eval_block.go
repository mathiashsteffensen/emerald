package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
)

func (vm *VM) EvalBlock(block object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
	bl := block.(*object.ClosedBlock)

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

	return vm.withExecutionContextForBlock(func() object.EmeraldValue {
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
	})
}

func (vm *VM) withExecutionContextForBlock(cb func() object.EmeraldValue) object.EmeraldValue {
	oldCtx := vm.ctx

	vm.ctx = vm.ctx.Outer

	val := cb()

	vm.ctx = oldCtx

	return val
}
