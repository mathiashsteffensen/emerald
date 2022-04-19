package core

import "emerald/object"

var BasicObject *object.Class

func init() {
	BasicObject = object.NewClass("BasicObject", nil, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
}
