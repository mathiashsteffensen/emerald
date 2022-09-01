package core

import "emerald/object"

var FalseClass *object.Class

var FALSE object.EmeraldValue

func InitFalseClass() {
	FalseClass = DefineClass(Object, "FalseClass", Object)

	DefineMethod(FalseClass, "to_s", falseToS())
	DefineMethod(FalseClass, "inspect", falseToS())

	FALSE = FalseClass.New()
}

func falseToS() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("false")
	}
}
