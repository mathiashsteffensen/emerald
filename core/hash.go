package core

import "emerald/object"

var Hash *object.Class

func InitHash() {
	Hash = object.NewClass("Hash", Object, Object.Class(), object.BuiltInMethodSet{
		"[]":   hashIndexAccessor(),
		"each": hashEach(),
	}, object.BuiltInMethodSet{}, Enumerable)
}

type HashInstance struct {
	*object.Instance
	Values map[string]object.EmeraldValue
	Keys   map[string]object.EmeraldValue
}

func NewHash() *HashInstance {
	return &HashInstance{
		Instance: Hash.New(),
		Values:   map[string]object.EmeraldValue{},
		Keys:     map[string]object.EmeraldValue{},
	}
}

func (hash *HashInstance) Get(key object.EmeraldValue) object.EmeraldValue {
	return hash.Values[key.HashKey()]
}

func (hash *HashInstance) Set(key object.EmeraldValue, value object.EmeraldValue) {
	hashKey := key.HashKey()
	hash.Values[hashKey] = value
	hash.Keys[hashKey] = key
}

func hashIndexAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		key := args[0]

		if value := target.(*HashInstance).Get(key); value != nil {
			return value
		}

		return NULL
	}
}

func hashEach() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		hash := target.(*HashInstance)

		for hashKey, value := range hash.Values {
			yield(block, hash.Keys[hashKey], value)
		}

		return hash
	}
}
