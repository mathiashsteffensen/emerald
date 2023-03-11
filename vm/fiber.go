package vm

import (
	"emerald/object"
	"fmt"
)

const (
	MaxStackSize = 2048
	MaxFrames    = 1024
)

var ErrStackOverflow = fmt.Errorf("stack overflow: max stack size of %d exceeded", MaxStackSize)

// Fiber is an abstract execution thread, separate from any OS level threads.
// Currently, this is kind of meaningless, but it is to allow for a concurrency implementation in the future.
type Fiber struct {
	// All fibers have their own stack allocated.
	// This allocates a []object.EmeraldValue of size MaxStackSize.
	// 32KB.
	stack []object.EmeraldValue
	// Always points to the next value. Top of stack is stack[sp-1]
	sp int
	// Call frames
	frames      []*Frame
	framesIndex int

	inRescue bool
}

func NewFiber(mainFrame *Frame) *Fiber {
	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &Fiber{
		stack:       make([]object.EmeraldValue, MaxStackSize),
		frames:      frames,
		framesIndex: 1,
	}
}

func (vm *VM) withFiber(cb func(fiber *Fiber)) {
	cb(vm.currentFiber())
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
		return nil
	}

	return fiber.stack[fiber.sp-1]
}

// push an obj on to the stack
func (fiber *Fiber) push(obj object.EmeraldValue) {
	if fiber.sp >= MaxStackSize {
		panic(ErrStackOverflow)
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
