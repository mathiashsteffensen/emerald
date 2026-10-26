package vm

import (
	"fmt"
)

// Backtrace describes the current fiber's call stack, innermost frame first, as Ruby does:
// "file:line:in 'method'". Native methods have no Frame; they are reported at their caller's line.
func (vm *VM) Backtrace() []string {
	fiber := vm.currentFiber()
	backtrace := make([]string, 0, fiber.framesIndex+len(fiber.nativeFrames))

	native := len(fiber.nativeFrames) - 1

	// Like Ruby, an error raised by Kernel#raise is reported from the line that called raise, not from raise itself
	if native >= 0 && fiber.nativeFrames[native].name == "raise" && fiber.nativeFrames[native].depth == fiber.framesIndex {
		native--
	}

	for i := fiber.framesIndex - 1; i >= 0; i-- {
		frame := fiber.frames[i]
		file, line := vm.frameSource(frame)

		// Natives called from frame i (depth i+1) run above it. Those called from deeper frames come first.
		for ; native >= 0 && fiber.nativeFrames[native].depth > i; native-- {
			if name := fiber.nativeFrames[native].name; name != "" {
				backtrace = append(backtrace, fmt.Sprintf("%s:%d:in '%s'", file, line, name))
			}
		}

		backtrace = append(backtrace, fmt.Sprintf("%s:%d:in '%s'", file, line, vm.frameName(frame, i == 0)))
	}

	return backtrace
}

func (vm *VM) frameSource(frame *Frame) (string, int) {
	block := frame.block

	file := block.FileName()
	if file == "" {
		file = block.File
	}
	if file == "" {
		file = vm.ctx.File
	}

	return file, block.LineAt(frame.ip)
}

func (vm *VM) frameName(frame *Frame, main bool) string {
	if name := frame.block.Name; name != "" {
		return name
	}

	if main {
		return "<main>"
	}

	return "<top (required)>"
}
