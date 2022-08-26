package core

import (
	"emerald/object"
)

var BasicObject *object.Class

func InitBasicObject() {
	BasicObject = object.NewClass("BasicObject", nil, nil, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
}
