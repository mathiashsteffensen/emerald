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

func (fiber *Fiber) currentFrame() *Frame {
	return fiber.frames[fiber.framesIndex-1]
}

// Adds a new call frame to the VM
func (fiber *Fiber) pushFrame(f *Frame) {
	fiber.frames[fiber.framesIndex] = f
	fiber.framesIndex++
}

func (fiber *Fiber) popFrame() *Frame {
	fiber.framesIndex--
	return fiber.frames[fiber.framesIndex]
}
