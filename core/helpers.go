package core

import (
	"emerald/object"
)

// Send is a function for calling methods that is dependency injected by the emerald/vm package
var Send func(self object.EmeraldValue, name string, block object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue

func DefineClass(namespace object.EmeraldValue, name string, super *object.Class) *object.Class {
	var superClass object.EmeraldValue

	if super != nil {
		superClass = super.Class()
	}

	class := object.NewClass(name, super, superClass, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	namespace.NamespaceDefinitionSet(name, class)
	class.SetParentNamespace(namespace)

	return class
}

func DefineModule(namespace object.EmeraldValue, name string) *object.Module {
	module := object.NewModule(name, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	if namespace != nil {
		namespace.NamespaceDefinitionSet(name, module)
		module.SetParentNamespace(namespace)
	}

	return module
}

func DefineMethod(receiver object.EmeraldValue, name string, method object.BuiltInMethod) {
	receiver.BuiltInMethodSet()[name] = method
}

func DefineSingletonMethod(receiver object.EmeraldValue, name string, method object.BuiltInMethod) {
	receiver.Class().BuiltInMethodSet()[name] = method
}

func NativeBoolToBooleanObject(input bool) object.EmeraldValue {
	if input {
		return TRUE
	}
	return FALSE
}

func IsError(obj object.EmeraldValue) bool {
	_, ok := obj.(object.EmeraldError)

	return ok
}
