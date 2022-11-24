package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/heap"
	"emerald/object"
	"fmt"
	"unicode"
)

func (vm *VM) executeOpConstantGet(ins compiler.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := heap.GetConstant(nameIndex).(*core.SymbolInstance).Value

	value, err := getConst(vm.ctx.Self, name)
	if err != nil {
		panic(err)
	}

	vm.push(value)
}

func (vm *VM) executeOpConstantSet(ins compiler.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := heap.GetConstant(nameIndex).(*core.SymbolInstance).Value
	// Don't pop it from the stack, we leave it there since assignment expressions return the assigned value
	value := vm.StackTop()

	setConst(vm.ctx.Self, name, value)
}

func (vm *VM) executeOpScopedConstantGet(ins compiler.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := heap.GetConstant(nameIndex).(*core.SymbolInstance).Value

	self := vm.pop()

	var result object.EmeraldValue

	if unicode.IsUpper(rune(name[0])) {
		value, err := getConst(self, name)
		if err != nil {
			err := core.NewStandardError(err.Error())
			core.Raise(err)
			return
		}
		result = value
	} else {
		result = vm.Send(self, name, core.NULL)
	}

	vm.push(result)
}

func getConst(self object.EmeraldValue, name string) (object.EmeraldValue, error) {
	value := self.NamespaceDefinitionGet(name)
	if value != nil {
		return value, nil
	}

	switch self.Type() {
	case object.INSTANCE_VALUE:
		// If it's an instance, check the class namespace
		value = self.Class().Super().NamespaceDefinitionGet(name)
	case object.STATIC_CLASS_VALUE:
		// If it's a singleton class, check the class namespace
		value = self.(*object.SingletonClass).Instance.NamespaceDefinitionGet(name)
	}

	if value == nil {
		// Try MainObject & Object as a last resort
		value = core.MainObject.NamespaceDefinitionGet(name)
		if value != nil {
			return value, nil
		}

		value = core.Object.NamespaceDefinitionGet(name)
		if value != nil {
			return value, nil
		}

		return nil, fmt.Errorf("uninitialized constant %s", name)

	}

	return value, nil
}

func setConst(self object.EmeraldValue, name string, value object.EmeraldValue) {
	self.NamespaceDefinitionSet(name, value)
	value.SetParentNamespace(self)
}
