package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"strings"
)

func (vm *VM) executeOpStringJoin(ins compiler.Instructions, ip int) {
	numStrings := int(vm.readUint8(ins, ip))

	stackPointer := vm.currentFiber().sp
	startPointer := stackPointer - numStrings

	objects := vm.stack()[startPointer:stackPointer]

	vm.currentFiber().sp = startPointer

	var out strings.Builder

	for _, o := range objects {
		stringified := vm.Send(o, "to_s", core.NULL, map[string]object.EmeraldValue{}).Inspect()

		out.WriteString(stringified)
	}

	vm.push(core.NewString(out.String()))
}
