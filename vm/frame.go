package vm

import (
	"emerald/bytecode"
	"emerald/object"
)

type Frame struct {
	block       *object.ClosedBlock
	ip          int
	basePointer int
	locals      map[int]*object.EmeraldValue
}

func NewFrame(block *object.ClosedBlock, basePointer int) *Frame {
	return &Frame{block: block, ip: -1, basePointer: basePointer}
}

func (f *Frame) Instructions() bytecode.Instructions {
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
	frame := fiber.frames[fiber.framesIndex]
	fiber.sp = frame.basePointer - 2
	if fiber.sp < 0 {
		fiber.sp = 0
	}
	return frame
}
