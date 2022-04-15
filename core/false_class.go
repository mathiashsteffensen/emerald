package core

import "emerald/object"

var FalseClass *object.Class

var FALSE object.EmeraldValue

func init() {
	FalseClass = object.NewClass("FalseClass", Object, object.BuiltInMethodSet{
		"==": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return nativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return nativeBoolToBooleanObject(target != args[0])
		},
	}, object.BuiltInMethodSet{})

	FALSE = FalseClass.New()
}
