package core

import "emerald/object"

var TrueClass *object.Class

var TRUE object.EmeraldValue

func InitTrueClass() {
	TrueClass = object.NewClass("TrueClass", Object, Object.Class(), object.BuiltInMethodSet{
		"==": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NativeBoolToBooleanObject(target != args[0])
		},
	}, object.BuiltInMethodSet{})

	TRUE = TrueClass.New()
	TRUE.Class().BuiltInMethodSet()["to_s"] = func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("true")
	}
}

func IsTruthy(obj object.EmeraldValue) bool {
	switch obj {
	case FALSE, NULL:
		return false
	default:
		return true
	}
}
