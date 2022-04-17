package vm

import (
	"emerald/compiler"
	"emerald/object"
)

type Frame struct {
	block       *object.ClosedBlock
	ip          int
	basePointer int
}

func NewFrame(block *object.ClosedBlock, basePointer int) *Frame {
	return &Frame{block: block, ip: -1, basePointer: basePointer}
}
func (f *Frame) Instructions() compiler.Instructions {
	return f.block.Instructions
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}
