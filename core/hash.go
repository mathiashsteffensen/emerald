package core

import "emerald/object"

var Hash = object.NewClass("Hash", Object, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

type HashInstance struct {
	*object.Instance
	Value map[object.EmeraldValue]object.EmeraldValue
}

func NewHash(val map[object.EmeraldValue]object.EmeraldValue) *HashInstance {
	return &HashInstance{
		Instance: Hash.New(),
		Value:    val,
	}
}
