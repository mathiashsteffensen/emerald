package core

import "emerald/object"

var NilClass *object.Class

var NULL object.EmeraldValue

func init() {
	NilClass = object.NewClass("NilClass", Object, object.BuiltInMethodSet{
		"==": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return nativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return nativeBoolToBooleanObject(target != args[0])
		},
	}, object.BuiltInMethodSet{})
	NULL = NilClass.New()
}
