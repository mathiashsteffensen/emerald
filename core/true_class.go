package core

import "emerald/object"

var TrueClass *object.Class

var TRUE object.EmeraldValue

func InitTrueClass() {
	TrueClass = DefineClass(Object, "TrueClass", Object)

	DefineMethod(TrueClass, "to_s", trueToS())
	DefineMethod(TrueClass, "inspect", trueToS())

	TRUE = TrueClass.New()
}

func trueToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("true")
	}
}

func IsTruthy(obj object.EmeraldValue) bool {
	switch obj {
	case FALSE, NULL:
		return false
	default:
		return true
	}
}
