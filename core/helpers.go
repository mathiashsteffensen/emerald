package core

import (
	"emerald/heap"
	"emerald/object"
	"fmt"
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

func EnforceArity(args []object.EmeraldValue, minArgs int, maxArgs int) ([]object.EmeraldValue, object.EmeraldError) {
	argsWithoutNilPointers := []object.EmeraldValue{}
	for _, arg := range args {
		if arg != nil {
			argsWithoutNilPointers = append(argsWithoutNilPointers, arg)
		}
	}
	numArgsGiven := len(argsWithoutNilPointers)

	if numArgsGiven < minArgs || numArgsGiven > maxArgs {
		err := NewArgumentError(numArgsGiven, maxArgs)
		Raise(err)
		return argsWithoutNilPointers, err
	}

	return argsWithoutNilPointers, nil
}

func EnforceArgumentType(typ *object.Class, arg object.EmeraldValue) object.EmeraldValue {
	argClass := arg.Class().Super().(*object.Class)
	if argClass.Name != typ.Name {
		err := NewTypeError(fmt.Sprintf("no implicit conversion of %s into %s", argClass.Name, typ.Name))
		Raise(err)
		return err
	}

	return nil
}

func Raise(err object.EmeraldError) {
	heap.SetGlobalVariableString("$!", err)
}

func NativeBoolToBooleanObject(input bool) object.EmeraldValue {
	if input {
		return TRUE
	}
	return FALSE
}
