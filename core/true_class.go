package core

import "emerald/object"

var TrueClass *object.Class

var TRUE object.EmeraldValue

func init() {
	TrueClass = object.NewClass("TrueClass", Object, object.BuiltInMethodSet{
		"==": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NativeBoolToBooleanObject(target != args[0])
		},
	}, object.BuiltInMethodSet{})

	TRUE = TrueClass.New()
}
