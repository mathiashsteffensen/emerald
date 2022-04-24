package core

import (
	"emerald/object"
)

var BasicObject *object.Class

func InitBasicObject() {
	BasicObject = object.NewClass("BasicObject", nil, Class, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
}
