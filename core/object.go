package core

import (
	"emerald/object"
)

var Object *object.Class

var MainObject *object.Instance

func InitObject() {
	Object = object.NewClass("Object", BasicObject, BasicObject.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	Object.Include(Kernel)

	DefineMethod(Object, "to_s", objectToS())
	DefineMethod(Object, "!@", objectBooleanNegate())
	DefineMethod(Object, "==", objectEquals())
	DefineMethod(Object, "!=", objectNotEquals())
	DefineMethod(Object, "methods", objectMethods())

	Object.NamespaceDefinitionSet(Object.Name, Object)
	Object.NamespaceDefinitionSet(Class.Name, Class)
	Class.SetParentNamespace(Object)
	Object.NamespaceDefinitionSet(Kernel.Name, Kernel)
	Kernel.SetParentNamespace(Object)

	MainObject = Object.New()

	DefineSingletonMethod(MainObject, "to_s", mainObjectToS())
	DefineSingletonMethod(MainObject, "inspect", mainObjectToS())
}

func objectToS() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(ctx.Self.Inspect())
	}
}

func mainObjectToS() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("main")
	}
}

func objectBooleanNegate() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return FALSE
	}
}

func objectEquals() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NativeBoolToBooleanObject(ctx.Self == args[0])
	}
}

func objectNotEquals() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NativeBoolToBooleanObject(ctx.Self != args[0])
	}
}

func objectMethods() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		methods := []object.EmeraldValue{}

		ancestors := ctx.Self.Ancestors()

		for _, ancestor := range ancestors {
			for _, method := range ancestor.Methods() {
				methods = append(methods, NewSymbol(method))
			}
		}

		return NewArray(methods)
	}
}
