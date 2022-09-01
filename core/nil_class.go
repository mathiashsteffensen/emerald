package core

import (
	"emerald/object"
	"reflect"
)

var NilClass *object.Class

var NULL *object.Instance

func InitNilClass() {
	NilClass = DefineClass(Object, "NilClass", Object)

	NULL = NilClass.New()

	DefineSingletonMethod(NULL, "to_s", nilToS())
	DefineSingletonMethod(NULL, "inspect", nilInspect())
}

func nilToS() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("")
	}
}

func nilInspect() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("nil")
	}
}

func IsNull(obj object.EmeraldValue) bool {
	if reflect.ValueOf(obj).IsNil() {
		return false
	}

	if obj == NULL {
		return true
	}

	return false
}
