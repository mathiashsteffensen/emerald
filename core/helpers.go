package core

import (
	"emerald/object"
)

// rt.Send is a function for calling methods that is dependency injected by the emerald/vm package

func (rt *Runtime) DefineClass(name string, super *object.Class) *object.Class {
	var superClass object.EmeraldValue

	if rt.Class != nil {
		superClass = rt.Class
	} else if super != nil {
		superClass = super.Class()
	}

	class := object.NewClass(name, super, superClass, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
	rt.Object.NamespaceDefinitionSet(name, class)
	class.SetParentNamespace(rt.Object)

	return class
}

func (rt *Runtime) DefineNestedClass(namespace object.EmeraldValue, name string, super *object.Class) *object.Class {
	var superClass object.EmeraldValue

	if rt.Class != nil {
		superClass = rt.Class
	} else if super != nil {
		superClass = super.Class()
	}

	class := object.NewClass(name, super, superClass, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
	namespace.NamespaceDefinitionSet(name, class)
	class.SetParentNamespace(namespace)

	return class
}

func (rt *Runtime) DefineModule(name string) *object.Module {
	module := object.NewModule(name, rt.Module, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
	rt.Object.NamespaceDefinitionSet(name, module)
	module.SetParentNamespace(rt.Object)

	return module
}

func (rt *Runtime) DefineNestedModule(namespace object.EmeraldValue, name string) *object.Module {
	module := object.NewModule(name, rt.Module, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
	namespace.NamespaceDefinitionSet(name, module)
	module.SetParentNamespace(namespace)

	return module
}

func (rt *Runtime) DefineMethod(receiver object.EmeraldValue, name string, method object.BuiltInMethod, visibilities ...object.MethodVisibility) {
	receiver.BuiltInMethodSet()[name] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}, Visibility: rt.getVisibility(visibilities)}
}

func (rt *Runtime) DefineSingletonMethod(receiver object.EmeraldValue, name string, method object.BuiltInMethod, visibilities ...object.MethodVisibility) {
	receiver.Class().BuiltInMethodSet()[name] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}, Visibility: rt.getVisibility(visibilities)}
}

func (rt *Runtime) getVisibility(visibilities []object.MethodVisibility) object.MethodVisibility {
	var visibility object.MethodVisibility
	if len(visibilities) != 0 {
		visibility = visibilities[0]
	} else {
		visibility = object.PUBLIC
	}

	return visibility
}

func (rt *Runtime) EnforceArity(
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
		err = rt.NewArgumentError(numArgsGiven, maxArgs)
		rt.Raise(err)
		return argsWithoutNilPointers, err
	}

	for _, kwarg := range requiredKwargs {
		if _, ok := kwargs[":"+kwarg]; !ok {
			err = rt.NewKeywordMissingArgumentError(kwarg)
			rt.Raise(err)
			return argsWithoutNilPointers, err
		}
	}

	return argsWithoutNilPointers, nil
}

func EnforceArgumentType[T object.EmeraldValue](rt *Runtime, typ *object.Class, arg object.EmeraldValue) (T, object.EmeraldError) {
	argClass := arg.Class().Super().(*object.Class)
	if argClass.Name != typ.Name {
		err := rt.NewNoConversionTypeError(typ.Name, argClass.Name)
		rt.Raise(err)
		var empty T
		return empty, err
	}

	return arg.(T), nil
}

func (rt *Runtime) Raise(err object.EmeraldError) object.EmeraldError {
	rt.Heap.SetGlobalVariableString("$!", err)
	if rt.OnRaise != nil {
		rt.OnRaise(err)
	}
	return err
}

func (rt *Runtime) RaiseGoError(err error) object.EmeraldError {
	emeraldErr := rt.NewStandardError(err.Error())
	return rt.Raise(emeraldErr)
}

func (rt *Runtime) NativeBoolToBooleanObject(input bool) object.EmeraldValue {
	if input {
		return rt.TRUE
	}
	return rt.FALSE
}
