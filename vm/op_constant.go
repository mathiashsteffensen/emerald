package vm

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/kernel"
	"emerald/object"
	"fmt"
)

func (vm *VM) executeOpConstantGet(ins compiler.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := kernel.GetConst(nameIndex).(*core.SymbolInstance).Value

	value, err := getConst(vm.ctx.Self, name)
	if err != nil {
		panic(err)
	}

	vm.push(value)
}

func (vm *VM) executeOpConstantSet(ins compiler.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := kernel.GetConst(nameIndex).(*core.SymbolInstance).Value
	// Don't pop it from the stack, we leave it there since assignment expressions return the assigned value
	value := vm.StackTop()

	setConst(vm.ctx.Self, name, value)
}

func (vm *VM) executeOpConstantGetOrSet(ins compiler.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := kernel.GetConst(nameIndex).(*core.SymbolInstance).Value
	valueIndex := compiler.ReadUint16(ins[ip+3:])
	vm.currentFrame().ip += 2

	self := vm.ctx.Self

	value, err := getConst(self, name)
	if err != nil {
		value = kernel.GetConst(valueIndex)

		// If self is an instance, define the constant on the class
		if self.Type() == object.INSTANCE_VALUE {
			self = self.Class().Super()
		}

		setConst(self, name, value)
	}

	vm.push(value)
}

func (vm *VM) executeOpScopedConstantGet(ins compiler.Instructions, ip int) {
	nameIndex := vm.readUint16(ins, ip)
	name := kernel.GetConst(nameIndex).(*core.SymbolInstance).Value

	self := vm.pop()

	value, err := getConst(self, name)
	if err != nil {
		panic(err)
	}

	vm.push(value)
}

func getConst(self object.EmeraldValue, name string) (object.EmeraldValue, error) {
	value := self.NamespaceDefinitionGet(name)
	if value == nil {
		// If it's an instance, check the class namespace
		if self.Type() == object.INSTANCE_VALUE {
			value = self.Class().Super().NamespaceDefinitionGet(name)
		}

		// If it's a singleton class, check the class namespace
		if self.Type() == object.STATIC_CLASS_VALUE {
			value = self.(*object.SingletonClass).Instance.NamespaceDefinitionGet(name)
		}

		if value == nil {
			// Try MainObject as a last resort
			value = core.MainObject.NamespaceDefinitionGet(name)

			if value == nil {
				return nil, fmt.Errorf("uninitialized constant %s", name)
			}
		}
	}
	return value, nil
}

func setConst(self object.EmeraldValue, name string, value object.EmeraldValue) {
	self.NamespaceDefinitionSet(name, value)
	value.SetParentNamespace(self)
}
