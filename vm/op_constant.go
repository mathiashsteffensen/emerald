package vm

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/debug"
	"emerald/object"
	"fmt"
	"unicode"
)

func (vm *VM) executeOpConstantGet(ins bytecode.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := vm.rt.Heap.GetConstant(nameIndex).Heap.(*core.SymbolInstance).Value

	value, err := getConst(vm.ctx.Self, name, vm.rt)
	if err != nil {
		debug.WarnF("\n%s", vm.currentFiber().currentFrame().block.Bytecode.InstructionSnapshot(ip))
		vm.rt.Raise(vm.rt.NewNameError(err.Error()))
		return
	}

	if !vm.ExceptionIsRaised() {
		vm.push(value)
	}
}

func (vm *VM) executeOpConstantSet(ins bytecode.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := vm.rt.Heap.GetConstant(nameIndex).Heap.(*core.SymbolInstance).Value
	// Don't pop it from the stack, we leave it there since assignment expressions return the assigned value
	value := vm.StackTop()

	setConst(vm.ctx.Self, name, value)
}

func (vm *VM) executeOpScopedConstantGet(ins bytecode.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := vm.rt.Heap.GetConstant(nameIndex).Heap.(*core.SymbolInstance).Value

	self := vm.pop()

	var result object.EmeraldValue

	if unicode.IsUpper(rune(name[0])) {
		value, err := getConst(self, name, vm.rt)
		if err != nil {
			err := vm.rt.NewStandardError(err.Error())
			vm.rt.Raise(err)
			return
		}
		result = value
	} else {
		result = vm.Send(self, name, vm.rt.NULL, map[string]object.EmeraldValue{})
	}

	if !vm.ExceptionIsRaised() {
		vm.push(result)
	}
}

func getConst(self object.EmeraldValue, name string, rt *core.Runtime) (object.EmeraldValue, error) {
	value := self.NamespaceDefinitionGet(name)
	if !value.IsNil() {
		return value, nil
	}

	switch self.Type() {
	case object.INSTANCE_VALUE:
		// If it's an instance, check the class namespace
		value = object.RealClass(self).NamespaceDefinitionGet(name)
	case object.STATIC_CLASS_VALUE:
		// If it's a singleton class, check the class namespace
		value = self.Heap.(*object.SingletonClass).Instance.NamespaceDefinitionGet(name)
	}

	if value.IsNil() {
		// Try MainObject & Object as a last resort
		value = rt.MainObject.NamespaceDefinitionGet(name)
		if !value.IsNil() {
			return value, nil
		}

		value = rt.Object.NamespaceDefinitionGet(name)
		if !value.IsNil() {
			return value, nil
		}

		return object.EmeraldValue{}, fmt.Errorf("uninitialized constant %s", name)
	}

	return value, nil
}

func setConst(self object.EmeraldValue, name string, value object.EmeraldValue) {
	self.NamespaceDefinitionSet(name, value)
	value.SetParentNamespace(self)
}
