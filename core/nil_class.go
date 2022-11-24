package core

import (
	"emerald/object"
)

var NilClass *object.Class

var NULL *object.Instance

func InitNilClass() {
	NilClass = DefineClass("NilClass", Object)

	NULL = NilClass.New()

	DefineSingletonMethod(NULL, "to_s", nilToS())
	DefineSingletonMethod(NULL, "inspect", nilInspect())
	DefineSingletonMethod(NULL, "!@", nilBooleanNegate())
}

func nilToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("")
	}
}

func nilInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("nil")
	}
}

func nilBooleanNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return TRUE
	}
}
