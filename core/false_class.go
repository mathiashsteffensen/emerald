package core

import "emerald/object"

var FalseClass *object.Class

var FALSE object.EmeraldValue

func InitFalseClass() {
	FalseClass = DefineClass("FalseClass", Object)

	DefineMethod(FalseClass, "to_s", falseToS())
	DefineMethod(FalseClass, "inspect", falseToS())
	DefineMethod(FalseClass, "!@", falseBooleanNegate())

	FALSE = FalseClass.New()
}

func falseToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("false")
	}
}

func falseBooleanNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return TRUE
	}
}
