package core

import (
	"emerald/object"
	"fmt"
	"math"
)

func (rt *Runtime) DefineClass(name string, super object.EmeraldValue) object.EmeraldValue {
	var superClass object.EmeraldValue

	if !rt.Class.IsNil() {
		superClass = rt.Class
	} else if !super.IsNil() {
		superClass = super.Class()
	}

	var superPtr *object.Class
	if !super.IsNil() {
		superPtr = super.Heap.(*object.Class)
	}

	class := object.NewClass(name, superPtr, superClass, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
	classVal := object.NewHeapObject(class)
	if !rt.Object.IsNil() {
		rt.Object.NamespaceDefinitionSet(name, classVal)
		class.SetParentNamespace(rt.Object)
	}

	return classVal
}

func (rt *Runtime) DefineNestedClass(namespace object.EmeraldValue, name string, super object.EmeraldValue) object.EmeraldValue {
	classVal := rt.DefineClass(name, super)
	namespace.NamespaceDefinitionSet(name, classVal)
	classVal.SetParentNamespace(namespace)
	return classVal
}

func (rt *Runtime) DefineModule(name string) object.EmeraldValue {
	module := object.NewModule(name, rt.Module, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
	moduleVal := object.NewHeapObject(module)
	if !rt.Object.IsNil() {
		rt.Object.NamespaceDefinitionSet(name, moduleVal)
		module.SetParentNamespace(rt.Object)
	}

	return moduleVal
}

func (rt *Runtime) DefineNestedModule(namespace object.EmeraldValue, name string) object.EmeraldValue {
	moduleVal := rt.DefineModule(name)
	namespace.NamespaceDefinitionSet(name, moduleVal)
	moduleVal.SetParentNamespace(namespace)

	return moduleVal
}

func (rt *Runtime) DefineMethod(receiver object.EmeraldValue, name string, method object.BuiltInMethod, visibilities ...object.MethodVisibility) {
	receiver.BuiltInMethodSet()[name] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}, Visibility: rt.getVisibility(visibilities)}
}

func (rt *Runtime) DefineSingletonMethod(receiver object.EmeraldValue, name string, method object.BuiltInMethod, visibilities ...object.MethodVisibility) {
	receiver.SingletonClass().BuiltInMethodSet()[name] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}, Visibility: rt.getVisibility(visibilities)}
}

func (rt *Runtime) DefineConstant(name string, value, namespace object.EmeraldValue) {
	namespace.NamespaceDefinitionSet(name, value)
}

func (rt *Runtime) DefineGlobalConstant(name string, value object.EmeraldValue) {
	rt.DefineConstant(name, value, rt.Object)
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

	numArgsGiven := len(args)

	if numArgsGiven < minArgs || numArgsGiven > maxArgs {
		err = rt.NewArgumentError(numArgsGiven, maxArgs)
		rt.Raise(err)
		return args, err
	}

	for _, kwarg := range requiredKwargs {
		if _, ok := kwargs[":"+kwarg]; !ok {
			err = rt.NewKeywordMissingArgumentError(kwarg)
			rt.Raise(err)
			return args, err
		}
	}

	return args, nil
}

func (rt *Runtime) EnforceIntegerArg(arg object.EmeraldValue) (int64, object.EmeraldError) {
	if !arg.Is(object.INTEGER_VALUE) {
		err := rt.NewNoConversionTypeError(rt.Integer.Heap.(*object.Class).Name, object.RealClass(arg).Heap.(*object.Class).Name)
		rt.Raise(err)
		return 0, err
	}
	return int64(arg.Num), nil
}

func (rt *Runtime) EnforceFloatArg(arg object.EmeraldValue) (float64, object.EmeraldError) {
	if !arg.Is(object.FLOAT_VALUE) {
		err := rt.NewNoConversionTypeError(rt.Float.Heap.(*object.Class).Name, object.RealClass(arg).Heap.(*object.Class).Name)
		rt.Raise(err)
		return 0, err
	}
	return math.Float64frombits(arg.Num), nil
}

func EnforceArgumentType[T any](rt *Runtime, typVal object.EmeraldValue, arg object.EmeraldValue) (T, object.EmeraldError) {
	typ := typVal.Heap.(*object.Class)
	realClass := object.RealClass(arg)
	argClass := realClass.Heap.(*object.Class)

	if argClass.Name != typ.Name {
		err := rt.NewNoConversionTypeError(typ.Name, argClass.Name)
		rt.Raise(err)
		var empty T
		return empty, err
	}

	if heap, ok := arg.Heap.(T); ok {
		return heap, nil
	}

	var empty T
	return empty, rt.RaiseGoError(fmt.Errorf("could not cast %T to %T", arg.Heap, empty))
}

func (rt *Runtime) Raise(err object.EmeraldError) object.EmeraldError {
	rt.Heap.SetGlobalVariableString("$!", object.NewHeapObject(err))
	if rt.OnRaise != nil {
		rt.OnRaise(err)
	}
	return err
}

func (rt *Runtime) ExceptionIsRaised() bool {
	exception := rt.Heap.GetGlobalVariableString("$!")

	return !exception.IsNil()
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
