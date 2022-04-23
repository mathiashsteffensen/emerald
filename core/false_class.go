package core

import "emerald/object"

var FalseClass *object.Class

var FALSE object.EmeraldValue

func init() {
	FalseClass = object.NewClass("FalseClass", Object, object.BuiltInMethodSet{
		"==": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NativeBoolToBooleanObject(target != args[0])
		},
		"to_s": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NewString("false")
		},
	}, object.BuiltInMethodSet{})

	FALSE = FalseClass.New()
}
