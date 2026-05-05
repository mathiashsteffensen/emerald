package vm

import (
	"emerald/object"
	"fmt"
)

const (
	InitialStackSize = 512
	MaxStackSize = 2048
	MaxFrames    = 1024
)

var ErrStackOverflow = fmt.Errorf("stack overflow: max stack size of %d exceeded", MaxStackSize)

// Fiber is an abstract execution thread, separate from any OS level threads.
// Currently, this is kind of meaningless, but it is to allow for a concurrency implementation in the future.
type Fiber struct {
	// The stack for this fiber.
	// Grows dynamically up to MaxStackSize.
	stack []object.EmeraldValue
	// Always points to the next value. Top of stack is stack[sp-1]
	sp int
	// Call frames
	frames      []*Frame
	framesIndex int
}

func NewFiber(mainFrame *Frame) *Fiber {
	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &Fiber{
		stack:       make([]object.EmeraldValue, InitialStackSize),
		frames:      frames,
		framesIndex: 1,
	}
}

func (vm *VM) currentFiber() *Fiber {
	return vm.fibers[vm.fiberIndex]
}

func (vm *VM) stack() []object.EmeraldValue {
	return vm.currentFiber().stack
}

// StackTop fetches the object at the top of the stack
func (fiber *Fiber) StackTop() object.EmeraldValue {
	if fiber.sp == 0 {
		return object.EmeraldValue{}
	}

	return fiber.stack[fiber.sp-1]
}

// push an obj on to the stack
func (fiber *Fiber) push(obj object.EmeraldValue) {
	if fiber.sp >= MaxStackSize {
		panic(ErrStackOverflow)
	}

	if fiber.sp == len(fiber.stack) {
		newSize := len(fiber.stack) * 2
		if newSize > MaxStackSize {
			newSize = MaxStackSize
		}

		newStack := make([]object.EmeraldValue, newSize)
		copy(newStack, fiber.stack)
		fiber.stack = newStack
	}

	fiber.stack[fiber.sp] = obj
	fiber.sp++
}

// pop an obj from the top of the stack
func (fiber *Fiber) pop() object.EmeraldValue {
	o := fiber.StackTop()
	fiber.sp--
	return o
}
