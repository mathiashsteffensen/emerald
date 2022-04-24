package core

import "emerald/object"

var TypeError *object.Class

func InitTypeError() {
	TypeError = object.NewClass("TypeError", StandardError, StandardError.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
}
