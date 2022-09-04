package vm

import "emerald/object"

const (
	StackSize = 2048
	MaxFrames = 1024
)

// Fiber is an abstract execution thread, separate from any OS level threads.
// Currently, this is kind of meaningless, but it is to allow for a concurrency implementation in the future.
type Fiber struct {
	// All fibers have their own stack allocated.
	// This allocates a []object.EmeraldValue of size StackSize.
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
		stack:       make([]object.EmeraldValue, StackSize),
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
