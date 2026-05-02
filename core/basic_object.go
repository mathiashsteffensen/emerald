package core

import (
	"emerald/object"
)

func (rt *Runtime) InitBasicObject() {
	rt.BasicObject = object.NewClass("BasicObject", nil, nil, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
}
