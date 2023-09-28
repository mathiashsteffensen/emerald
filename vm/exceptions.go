package vm

import (
	"emerald/core"
	"emerald/debug"
	"emerald/heap"
	"emerald/object"
)

func (vm *VM) ExceptionIsRaised() bool {
	globalException := heap.GetGlobalVariableString("$!")

	if globalException == nil && globalException != core.NULL {
		return false
	}

	return true
}

func (vm *VM) popFramesUntilExceptionRescuedOrProgramTerminates() bool {
	// If this is called we assume that this is not nil
	raisedException := heap.GetGlobalVariableString("$!").(object.EmeraldError)

	debug.InternalDebugF("Raising %s in frame %d", raisedException.ClassName(), vm.currentFiber().framesIndex)

	rescued := true

	vm.withFiber(func(fiber *Fiber) {
		for !fiber.currentFrame().rescuesException(raisedException) {
			fiber.popFrame()
			if fiber.framesIndex == 0 {
				rescued = false
				break
			}
		}

		if rescued {
			debug.InternalDebugF("Rescued in frame %d!", fiber.framesIndex)
			vm.currentFiber().inRescue = true
			debug.InternalDebug("Evaluating rescue clause")
			vm.rawEvalBlock(fiber.currentFrame().blockRescuingException(raisedException), core.NULL, map[string]object.EmeraldValue{})
			heap.SetGlobalVariableString("$!", nil)
			debug.InternalDebugF("Done evaluating rescue clause, framesIndex=%d", fiber.framesIndex)
			vm.currentFiber().inRescue = false
		}
	})

	return rescued
}
