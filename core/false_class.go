package core

import "emerald/object"

func (rt *Runtime) InitFalseClass() {
	rt.FalseClass = rt.DefineClass("FalseClass", rt.Object)

	rt.DefineMethod(rt.FalseClass, "to_s", rt.falseToS())
	rt.DefineMethod(rt.FalseClass, "inspect", rt.falseToS())
	rt.DefineMethod(rt.FalseClass, "!@", rt.falseBooleanNegate())

	rt.FALSE = object.EmeraldValue{TypeID: object.FALSE_VALUE, Heap: rt.FalseClass.Heap}
}

func (rt *Runtime) falseToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString("false")
	}
}

func (rt *Runtime) falseBooleanNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.TRUE
	}
}
