package core

import "emerald/object"

var Numeric *object.Class

func InitNumeric() {
	Numeric = DefineClass(Object, "Numeric", Object)
}
