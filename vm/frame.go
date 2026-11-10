package vm

import (
	"emerald/bytecode"
	"emerald/object"
)

type Frame struct {
	block        *object.ClosedBlock
	ip           int
	basePointer  int
	locals       map[int]*object.EmeraldValue
	rescuedError object.EmeraldValue
	// argGiven marks which parameters received an argument for methods with
	// default parameters. It is nil when the method has no defaults.
	argGiven []bool
	// home is this method invocation, when block literals it creates can `return` from it
	home *object.Invocation
	// ensuring is the return or raise that the running ensure clause interrupted
	ensuring *interruptedExit
}

// interruptedExit is a return with returnValue, or a raise of err when err is set
type interruptedExit struct {
	returnValue object.EmeraldValue
	err         object.EmeraldError
}

func NewFrame(block *object.ClosedBlock, basePointer int) *Frame {
	return &Frame{block: block, ip: -1, basePointer: basePointer}
}

// homeInvocation is the method invocation a `return` in this frame's block leaves
func (f *Frame) homeInvocation() *object.Invocation {
	if f.home != nil {
		return f.home
	}

	return f.block.Home
}

// ensureEntry is the ensure clause that protects the current instruction, or nil
func (f *Frame) ensureEntry() *object.ExceptionTableEntry {
	for i, entry := range f.block.ExceptionTable {
		if entry.Ensure && f.ip >= entry.StartIP && f.ip <= entry.EndIP {
			return &f.block.ExceptionTable[i]
		}
	}

	return nil
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
