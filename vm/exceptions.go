package vm

import (
	"emerald/object"
)

// runLevel is one runWhile loop, executing frames from frameIndex and up
type runLevel struct {
	frameIndex int
}

// rescueSignal is panicked by handleRaise to unwind the Go stack, through any native methods, to the runWhile
// loop executing the frame with the matching rescue handler
type rescueSignal struct {
	level *runLevel
}

func (vm *VM) handleRaise(err object.EmeraldError) {
	fiber := vm.currentFiber()

	// Capture before unwinding. A re-raised error keeps the backtrace from where it was first raised.
	if err.Backtrace() == nil {
		err.SetBacktrace(vm.Backtrace())
	}

	for fiber.framesIndex > 0 {
		frame := fiber.currentFrame()
		block := frame.block

		for _, entry := range block.ExceptionTable {
			if frame.ip >= entry.StartIP && frame.ip <= entry.EndIP {
				matches := false

				if len(entry.CaughtErrorClasses) == 0 {
					class := vm.rt.StandardError
					if vm.rt.IsTruthy(vm.rt.Send(object.NewHeapObject(err), "is_a?", vm.rt.NULL, map[string]object.EmeraldValue{}, class)) {
						matches = true
					}
				} else {
					for _, className := range entry.CaughtErrorClasses {
						class := vm.rt.Object.NamespaceDefinitionGet(className)
						if !class.IsNil() && vm.rt.IsTruthy(vm.rt.Send(object.NewHeapObject(err), "is_a?", vm.rt.NULL, map[string]object.EmeraldValue{}, class)) {
							matches = true
							break
						}
					}
				}

				if matches {
					frame.ip = entry.HandlerIP - 1
					fiber.sp = frame.basePointer + block.NumLocals
					frame.rescuedError = object.NewHeapObject(err)
					vm.rt.Heap.SetGlobalVariableString("$!", object.EmeraldValue{})

					// Natives between here and the handler stop on $!, which is now cleared, so skip past them
					if level := vm.runLevelFor(fiber.framesIndex); level != nil {
						panic(rescueSignal{level: level})
					}
					return
				}
			}
		}

		fiber.popFrame()
	}
}

// runLevelFor finds the innermost runWhile loop executing the frame at frameIndex
func (vm *VM) runLevelFor(frameIndex int) *runLevel {
	for i := len(vm.runLevels) - 1; i >= 0; i-- {
		if vm.runLevels[i].frameIndex <= frameIndex {
			return vm.runLevels[i]
		}
	}

	return nil
}

func (vm *VM) ExceptionIsRaised() bool {
	return vm.rt.ExceptionIsRaised()
}
