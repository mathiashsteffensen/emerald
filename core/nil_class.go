package core

import (
	"emerald/object"
	"reflect"
)

var NilClass *object.Class

var NULL object.EmeraldValue

func init() {
	NilClass = object.NewClass("NilClass", Object, object.BuiltInMethodSet{
		"==": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NativeBoolToBooleanObject(target != args[0])
		},
	}, object.BuiltInMethodSet{})
	NULL = NilClass.New()
}

func IsNull(obj object.EmeraldValue) bool {
	if reflect.ValueOf(obj).IsNil() {
		return false
	}

	if super, ok := obj.(*object.Class); ok && super.Name == NilClass.Name {
		return true
	}

	return false
}
