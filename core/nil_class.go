package core

import (
	"emerald/object"
)

func (rt *Runtime) InitNilClass() {
	rt.NilClass = rt.DefineClass("NilClass", rt.Object)

	rt.NULL = object.EmeraldValue{TypeID: object.NIL_VALUE, Heap: rt.NilClass.Heap}

	rt.DefineMethod(rt.NilClass, "to_s", rt.nilToS())
	rt.DefineMethod(rt.NilClass, "inspect", rt.nilInspect())
	rt.DefineMethod(rt.NilClass, "!@", rt.nilBooleanNegate())
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
