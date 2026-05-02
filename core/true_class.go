package core

import "emerald/object"

func (rt *Runtime) InitTrueClass() {
	rt.TrueClass = rt.DefineClass("TrueClass", rt.Object)

	rt.DefineMethod(rt.TrueClass, "to_s", rt.trueToS())
	rt.DefineMethod(rt.TrueClass, "inspect", rt.trueToS())

	rt.TRUE = rt.TrueClass.New()
}

func (rt *Runtime) trueToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString("true")
	}
}

func (rt *Runtime) IsTruthy(obj object.EmeraldValue) bool {
	switch obj {
	case rt.FALSE, rt.NULL:
		return false
	default:
		return true
	}
}
