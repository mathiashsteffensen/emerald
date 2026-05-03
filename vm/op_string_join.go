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

	objects := vm.stack()[startPointer:stackPointer]

	vm.currentFiber().sp = startPointer

	var out strings.Builder

	for _, o := range objects {
		res := vm.Send(o, "to_s", vm.rt.NULL, map[string]object.EmeraldValue{})

		if vm.ExceptionIsRaised() {
			return
		}

		stringified := res.Inspect()

		out.WriteString(stringified)
	}

	vm.push(vm.rt.NewString(out.String()))
}
