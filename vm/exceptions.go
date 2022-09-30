package vm

import (
	"emerald/core"
	"emerald/heap"
	"emerald/log"
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

	log.InternalDebugF("Raising %s in frame %d", raisedException.Inspect(), vm.currentFiber().framesIndex)

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
			log.InternalDebug("Rescued!")
			vm.currentFiber().inRescue = true
			log.InternalDebug("Evaluating rescue clause")
			vm.rawEvalBlock(fiber.currentFrame().blockRescuingException(raisedException), core.NULL)
			heap.SetGlobalVariableString("$!", nil)
			log.InternalDebugF("Done evaluating rescue clause, framesIndex=%d", fiber.framesIndex)
			vm.currentFiber().inRescue = false
		}
	})

	return rescued
}
