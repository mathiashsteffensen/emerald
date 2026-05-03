package vm

import (
	"emerald/object"
)

func (vm *VM) handleRaise(err object.EmeraldError) {
	fiber := vm.currentFiber()

	for fiber.framesIndex > 0 {
		frame := fiber.currentFrame()
		block := frame.block

		for _, entry := range block.ExceptionTable {
			if frame.ip >= entry.StartIP && frame.ip <= entry.EndIP {
				matches := false

				if len(entry.CaughtErrorClasses) == 0 {
					class := vm.rt.StandardError
					if vm.rt.IsTruthy(vm.rt.Send(err, "is_a?", vm.rt.NULL, map[string]object.EmeraldValue{}, class)) {
						matches = true
					}
				} else {
					for _, className := range entry.CaughtErrorClasses {
						class := vm.rt.Object.NamespaceDefinitionGet(className)
						if class != nil && vm.rt.IsTruthy(vm.rt.Send(err, "is_a?", vm.rt.NULL, map[string]object.EmeraldValue{}, class)) {
							matches = true
							break
						}
					}
				}

				if matches {
					frame.ip = entry.HandlerIP - 1
					fiber.sp = frame.basePointer + block.NumLocals
					vm.rt.Heap.SetGlobalVariableString("$!", nil)
					return
				}
			}
		}

		fiber.popFrame()
	}
}

func (vm *VM) ExceptionIsRaised() bool {
	globalException := vm.rt.Heap.GetGlobalVariableString("$!")

	if globalException == nil || globalException == vm.rt.NULL {
		return false
	}

	return true
}
