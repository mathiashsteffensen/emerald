package vm

import "emerald/object"

// enterEnsure starts the ensure clause that protects the current instruction,
// to run before the frame returns returnValue. It is false when no clause protects it.
func (vm *VM) enterEnsure(returnValue object.EmeraldValue) bool {
	fiber := vm.currentFiber()
	frame := fiber.currentFrame()

	entry := frame.ensureEntry()
	if entry == nil {
		return false
	}

	frame.ensuring = &interruptedExit{returnValue: returnValue}
	frame.ip = entry.HandlerIP - 1
	fiber.sp = frame.basePointer + frame.block.NumLocals

	return true
}

// executeOpEnsureEnd finishes the return, or raises the error again, that the ensure clause interrupted
func (vm *VM) executeOpEnsureEnd() {
	frame := vm.currentFiber().currentFrame()
	exit := frame.ensuring
	frame.ensuring = nil

	if exit.err != nil {
		vm.rt.Raise(exit.err)
		return
	}

	vm.currentFiber().popFrame()

	vm.push(exit.returnValue)
}
