package vm

import (
	"emerald/core"
	"emerald/object"
)

func (vm *VM) newContext(file string, self, block object.EmeraldValue) *object.Context {
	return &object.Context{
		File:                    file,
		Self:                    self,
		Block:                   block,
		Yield:                   vm.Yield,
		BlockGiven:              vm.BlockGiven,
		DefaultMethodVisibility: object.PUBLIC,
	}
}

func (vm *VM) newEnclosedContext(file string, self, block object.EmeraldValue) *object.Context {
	ctx := vm.newContext(file, self, block)
	ctx.Outer = vm.ctx
	return ctx
}

func (vm *VM) BlockGiven() bool {
	return vm.ctx.Block != core.NULL
}
