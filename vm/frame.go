package vm

import (
	"emerald/compiler"
	"emerald/core"
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
	frame := fiber.frames[fiber.framesIndex]
	fiber.sp = frame.basePointer - 3
	if fiber.sp < 0 {
		fiber.sp = 0
	}
	return frame
}

func (f *Frame) blockRescuingException(exception object.EmeraldError) *object.ClosedBlock {
	for _, rescueBlock := range f.block.RescueBlocks {
		caughtClassName := rescueBlock.CaughtErrorClasses.Find(func(className string) bool {
			class := core.Object.NamespaceDefinitionGet(className)

			return core.IsTruthy(core.Send(exception, "is_a?", core.NULL, class))
		})

		if caughtClassName != nil {
			block := object.NewBlock(rescueBlock.Instructions, 0, 0, []string{})
			return object.NewClosedBlock(f.block.Context, block, []object.EmeraldValue{}, f.block.File, object.PUBLIC)
		}
	}

	return nil
}

func (f *Frame) rescuesException(exception object.EmeraldError) bool {
	return f.blockRescuingException(exception) != nil
}
