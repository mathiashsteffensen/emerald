package core

import "emerald/object"

var Hash *object.Class

func InitHash() {
	Hash = object.NewClass("Hash", Object, Object.Class(), object.BuiltInMethodSet{
		"[]": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			key := args[0]

			if value := target.(*HashInstance).Value[key.HashKey()]; value != nil {
				return value
			}

			return NULL
		},
	}, object.BuiltInMethodSet{})
}

type HashInstance struct {
	*object.Instance
	Value map[string]object.EmeraldValue
}

func NewHash(val map[string]object.EmeraldValue) *HashInstance {
	return &HashInstance{
		Instance: Hash.New(),
		Value:    val,
	}
}
