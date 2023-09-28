package core

import (
	"emerald/heap"
	"emerald/object"
)

// Send is a function for calling methods that is dependency injected by the emerald/vm package
var Send func(self object.EmeraldValue, name string, block object.EmeraldValue, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue

func DefineClass(name string, super *object.Class) *object.Class {
	var superClass object.EmeraldValue

	if super != nil {
		superClass = super.Class()
	}

	class := object.NewClass(name, super, superClass, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	Object.NamespaceDefinitionSet(name, class)
	class.SetParentNamespace(Object)

	return class
}

func DefineNestedClass(namespace object.EmeraldValue, name string, super *object.Class) *object.Class {
	var superClass object.EmeraldValue

	if super != nil {
		superClass = super.Class()
	}

	class := object.NewClass(name, super, superClass, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	namespace.NamespaceDefinitionSet(name, class)
	class.SetParentNamespace(namespace)

	return class
}

func DefineModule(name string) *object.Module {
	module := object.NewModule(name, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	Object.NamespaceDefinitionSet(name, module)
	module.SetParentNamespace(Object)

	return module
}

func DefineNestedModule(namespace object.EmeraldValue, name string) *object.Module {
	module := object.NewModule(name, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	Object.NamespaceDefinitionSet(name, module)
	module.SetParentNamespace(Object)

	return module
}

func DefineMethod(receiver object.EmeraldValue, name string, method object.BuiltInMethod, visibilities ...object.MethodVisibility) {
	receiver.BuiltInMethodSet()[name] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}, Visibility: getVisibility(visibilities)}
}

func DefineSingletonMethod(receiver object.EmeraldValue, name string, method object.BuiltInMethod, visibilities ...object.MethodVisibility) {
	receiver.Class().BuiltInMethodSet()[name] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}, Visibility: getVisibility(visibilities)}
}

func getVisibility(visibilities []object.MethodVisibility) object.MethodVisibility {
	var visibility object.MethodVisibility
	if len(visibilities) != 0 {
		visibility = visibilities[0]
	} else {
		visibility = object.PUBLIC
	}

	return visibility
}

func EnforceArity(
	args []object.EmeraldValue,
	kwargs map[string]object.EmeraldValue,
	minArgs int,
	maxArgs int,
	requiredKwargs ...string,
) ([]object.EmeraldValue, object.EmeraldError) {
	var err object.EmeraldError

	argsWithoutNilPointers := []object.EmeraldValue{}
	for _, arg := range args {
		if arg != nil {
			argsWithoutNilPointers = append(argsWithoutNilPointers, arg)
		}
	}
	numArgsGiven := len(argsWithoutNilPointers)

	if numArgsGiven < minArgs || numArgsGiven > maxArgs {
		err = NewArgumentError(numArgsGiven, maxArgs)
		Raise(err)
		return argsWithoutNilPointers, err
	}

	for _, kwarg := range requiredKwargs {
		if _, ok := kwargs[":"+kwarg]; !ok {
			err = NewKeywordMissingArgumentError(kwarg)
			Raise(err)
			return argsWithoutNilPointers, err
		}
	}

	return argsWithoutNilPointers, nil
}

func EnforceArgumentType[T object.EmeraldValue](typ *object.Class, arg object.EmeraldValue) (T, object.EmeraldError) {
	argClass := arg.Class().Super().(*object.Class)
	if argClass.Name != typ.Name {
		err := NewNoConversionTypeError(typ.Name, argClass.Name)
		Raise(err)
		var empty T
		return empty, err
	}

	return arg.(T), nil
}

func Raise(err object.EmeraldError) object.EmeraldError {
	heap.SetGlobalVariableString("$!", err)
	return err
}

func RaiseGoError(err error) object.EmeraldError {
	emeraldErr := NewStandardError(err.Error())
	return Raise(emeraldErr)
}

func NativeBoolToBooleanObject(input bool) object.EmeraldValue {
	if input {
		return TRUE
	}
	return FALSE
}
