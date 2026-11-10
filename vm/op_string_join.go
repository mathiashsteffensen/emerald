package vm

import (
	"emerald/bytecode"
	"emerald/object"
	"strings"
)

func (vm *VM) executeOpStringJoin(ins bytecode.Instructions, ip int) {
	numStrings := int(vm.readUint8(ins, ip))

	stackPointer := vm.currentFiber().sp
	startPointer := stackPointer - numStrings

	// The parts stay on the stack until the end, so a Ruby to_s cannot push over them
	objects := vm.stack()[startPointer:stackPointer]

	var out strings.Builder

	for _, o := range objects {
		res := vm.Send(o, "to_s", vm.rt.NULL, map[string]object.EmeraldValue{})

		if vm.ExceptionIsRaised() {
			return
		}

		stringified := res.Inspect()

		out.WriteString(stringified)
	}

	vm.currentFiber().sp = startPointer
	vm.push(vm.rt.NewString(out.String()))
}
