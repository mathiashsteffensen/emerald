package core

import (
	"emerald/object"
)

var Object *object.Class

var MainObject *object.Instance

func init() {
	Object = object.NewClass(
		"Object",
		BasicObject,
		object.BuiltInMethodSet{
			"methods": func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				methods := []object.EmeraldValue{}

				for key := range target.(*object.Instance).DefinedMethodSet() {
					methods = append(methods, NewString(key))
				}

				return NewArray(methods)
			},
			"to_s": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewString(target.Inspect())
			},
			"==": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NativeBoolToBooleanObject(target == args[0])
			},
			"!=": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NativeBoolToBooleanObject(target.Inspect() != args[0].Inspect())
			},
		},
		object.BuiltInMethodSet{
			"ancestors": func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewArray(target.Ancestors())
			},
			"new": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return target.(*object.Class).New()
			},
			"define_method": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				target.DefineMethod(false, block, args...)

				return args[0]
			},
			"to_s": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewString(target.Inspect())
			},
		},
		Kernel,
	)

	MainObject = Object.New()
	MainObject.BuiltInSingletonMethods["to_s"] = func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("main:Object")
	}
}
