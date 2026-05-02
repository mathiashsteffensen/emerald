package core

import (
	"emerald/object"
)

func (rt *Runtime) InitNilClass() {
	rt.NilClass = rt.DefineClass("NilClass", rt.Object)

	rt.NULL = rt.NilClass.New()

	rt.DefineSingletonMethod(rt.NULL, "to_s", rt.nilToS())
	rt.DefineSingletonMethod(rt.NULL, "inspect", rt.nilInspect())
	rt.DefineSingletonMethod(rt.NULL, "!@", rt.nilBooleanNegate())
}

func (rt *Runtime) nilToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString("")
	}
}

func (rt *Runtime) nilInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString("nil")
	}
}

func (rt *Runtime) nilBooleanNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.TRUE
	}
}
