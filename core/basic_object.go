package core

import (
	"emerald/object"
)

func (rt *Runtime) InitBasicObject() {
	rt.BasicObject = object.NewHeapObject(object.NewClass("BasicObject", nil, rt.Class, object.BuiltInMethodSet{}, object.BuiltInMethodSet{}))
}
