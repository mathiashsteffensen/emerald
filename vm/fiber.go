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
}

func NewFiber() *Fiber {
	return &Fiber{stack: make([]object.EmeraldValue, StackSize)}
}
