package vm

import (
	"emerald/bytecode"
	"emerald/object"
)

func (vm *VM) executeOpCheckCaseEqual(ins bytecode.Instructions, ip int) {
	numMatchers := int(vm.readUint8(ins, ip))
	jumpPositionIfNoMatch := int(vm.readUint16(ins, ip+1))

	// This is essentially popping the top numMatchers elements from the stack
	// But this way is faster than calling vm.pop() in a loop
	matchers := vm.stack()[vm.currentFiber().sp-numMatchers : vm.currentFiber().sp]
	vm.currentFiber().sp -= numMatchers

	subject := vm.StackTop()

	for _, matcher := range matchers {
		res := vm.Send(matcher, "===", vm.rt.NULL, map[string]object.EmeraldValue{}, subject)

		if vm.ExceptionIsRaised() {
			return
		}

		if res.Is(object.TRUE_VALUE) {
			vm.pop()
			return
		}
	}

	vm.currentFiber().currentFrame().ip = jumpPositionIfNoMatch - 1
}
