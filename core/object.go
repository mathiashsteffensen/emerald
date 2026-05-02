package core

import (
	"emerald/object"
)

func (rt *Runtime) InitObject() {
	rt.Object = object.NewClass("Object", rt.BasicObject, rt.BasicObject.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	rt.Object.Include(rt.Kernel)

	rt.DefineMethod(rt.Object, "to_s", rt.objectToS())
	rt.DefineMethod(rt.Object, "!@", rt.objectBooleanNegate())
	rt.DefineMethod(rt.Object, "==", rt.objectEquals())
	rt.DefineMethod(rt.Object, "!=", rt.objectNotEquals())
	rt.DefineMethod(rt.Object, "methods", rt.objectMethods())

	rt.Object.NamespaceDefinitionSet("Object", rt.Object)
	rt.Object.NamespaceDefinitionSet("Class", rt.Class)
	rt.Class.SetParentNamespace(rt.Object)
	rt.Object.NamespaceDefinitionSet("Kernel", rt.Kernel)
	rt.Kernel.SetParentNamespace(rt.Object)

	rt.MainObject = rt.Object.New()

	rt.DefineSingletonMethod(rt.MainObject, "to_s", rt.mainObjectToS())
	rt.DefineSingletonMethod(rt.MainObject, "inspect", rt.mainObjectToS())
}

func (rt *Runtime) objectToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(ctx.Self.Inspect())
	}
}

func (rt *Runtime) mainObjectToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString("main")
	}
}

func (rt *Runtime) objectBooleanNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.FALSE
	}
}

func (rt *Runtime) objectEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NativeBoolToBooleanObject(ctx.Self == args[0])
	}
}

func (rt *Runtime) objectNotEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NativeBoolToBooleanObject(ctx.Self != args[0])
	}
}

func (rt *Runtime) objectMethods() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		methods := []object.EmeraldValue{}

		ancestors := ctx.Self.Ancestors()

		for _, ancestor := range ancestors {
			for _, method := range ancestor.Methods() {
				methods = append(methods, rt.NewSymbol(method))
			}
		}

		return rt.NewArray(methods)
	}
}
