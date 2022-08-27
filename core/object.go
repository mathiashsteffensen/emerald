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
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(target.Inspect())
	}
}

func mainObjectToS() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("main")
	}
}

func objectEquals() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NativeBoolToBooleanObject(self == args[0])
	}
}

func objectNotEquals() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NativeBoolToBooleanObject(self != args[0])
	}
}

func objectMethods() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		methods := []object.EmeraldValue{}

		for _, method := range self.Methods(self) {
			methods = append(methods, NewSymbol(method))
		}

		return NewArray(methods)
	}
}
