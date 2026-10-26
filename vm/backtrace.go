package vm

import (
	"fmt"
)

// Backtrace describes the current fiber's call stack, innermost frame first, as Ruby does:
// "file:line:in 'method'". Native methods have no frame and do not appear.
func (vm *VM) Backtrace() []string {
	fiber := vm.currentFiber()
	backtrace := make([]string, 0, fiber.framesIndex)

	for i := fiber.framesIndex - 1; i >= 0; i-- {
		backtrace = append(backtrace, vm.frameLocation(fiber.frames[i], i == 0))
	}

	return backtrace
}

func (vm *VM) frameLocation(frame *Frame, main bool) string {
	block := frame.block

	file := block.FileName()
	if file == "" {
		file = block.File
	}
	if file == "" {
		file = vm.ctx.File
	}

	name := block.Name
	if name == "" {
		if main {
			name = "<main>"
		} else {
			name = "<top (required)>"
		}
	}

	return fmt.Sprintf("%s:%d:in '%s'", file, block.LineAt(frame.ip), name)
}
