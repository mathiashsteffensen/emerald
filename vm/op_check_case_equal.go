package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
)

func (vm *VM) executeOpCheckCaseEqual(ins compiler.Instructions, ip int) {
	numMatchers := int(vm.readUint8(ins, ip))
	jumpPositionIfNoMatch := int(vm.readUint16(ins, ip+1))

	// This is essentially popping the top numMatchers elements from the stack
	// But this way is faster than calling vm.pop() in a loop
	matchers := vm.stack()[vm.currentFiber().sp-numMatchers : vm.currentFiber().sp]
	vm.currentFiber().sp -= numMatchers

	subject := vm.StackTop()

	for _, matcher := range matchers {
		if vm.Send(matcher, "===", core.NULL, map[string]object.EmeraldValue{}, subject) == core.TRUE {
			vm.pop()
			return
		}
	}

	vm.currentFiber().currentFrame().ip = jumpPositionIfNoMatch - 1
}
