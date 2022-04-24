package core

import (
	"emerald/object"
)

var Object *object.Class

var MainObject *object.Instance

func InitObject() {
	Object = object.NewClass(
		"Object",
		BasicObject,
		BasicObject.Class(),
		object.BuiltInMethodSet{
			"to_s": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewString(target.Inspect())
			},
			"==": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NativeBoolToBooleanObject(target == args[0])
			},
			"!=": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NativeBoolToBooleanObject(target.Inspect() != args[0].Inspect())
			},
			"methods": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				methods := []object.EmeraldValue{}

				for _, method := range target.Methods(target) {
					methods = append(methods, NewSymbol(method))
				}

				return NewArray(methods)
			},
		},
		object.BuiltInMethodSet{
			"to_s": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewString(target.Inspect())
			},
		},
		Kernel,
	)

	MainObject = Object.New()
	MainObject.Class().BuiltInMethodSet()["to_s"] = func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("main:Object")
	}
}
