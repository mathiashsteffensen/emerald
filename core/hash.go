package core

import (
	"emerald/object"
)

var Hash *object.Class

func InitHash() {
	Hash = DefineClass(Object, "Hash", Object)

	Hash.Include(Enumerable)

	DefineMethod(Hash, "[]", hashIndexAccessor())
	DefineMethod(Hash, "==", hashEquals())
	DefineMethod(Hash, "each", hashEach())
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
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		key := args[0]

		if value := ctx.Self.(*HashInstance).Get(key); value != nil {
			return value
		}

		return NULL
	}
}

func hashEach() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		hash := ctx.Self.(*HashInstance)

		for hashKey, value := range hash.Values {
			ctx.Yield(hash.Keys[hashKey], value)
		}

		return hash
	}
}

func hashEquals() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		args, err := EnforceArity(args, 1, 1)

		if err != nil {
			return err
		}

		otherObj := args[0]
		if otherObj.Class().Super() != Hash {
			return FALSE
		}

		hash := ctx.Self.(*HashInstance)
		otherHash := otherObj.(*HashInstance)

		for key, value := range hash.Values {
			otherValue, ok := otherHash.Values[key]
			if !ok {
				return FALSE
			}

			if Send(value, "==", NULL, otherValue) != TRUE {
				return FALSE
			}
		}

		return TRUE
	}
}
