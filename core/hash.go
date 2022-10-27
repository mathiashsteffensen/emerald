package core

import (
	"emerald/object"
	"github.com/elliotchance/orderedmap/v2"
)

var Hash *object.Class

func InitHash() {
	Hash = DefineClass(Object, "Hash", Object)

	Hash.Include(Enumerable)

	DefineMethod(Hash, "[]", hashIndexAccessor())
	DefineMethod(Hash, "[]=", hashIndexSetter())
	DefineMethod(Hash, "==", hashEquals())
	DefineMethod(Hash, "each", hashEach())
}

type HashInstance struct {
	*object.Instance
	Values *orderedmap.OrderedMap[string, object.EmeraldValue] // Only Values need to be ordered since we always iterate on those
	Keys   map[string]object.EmeraldValue
}

func NewHash() *HashInstance {
	return &HashInstance{
		Instance: Hash.New(),
		Values:   orderedmap.NewOrderedMap[string, object.EmeraldValue](),
		Keys:     map[string]object.EmeraldValue{},
	}
}

func (hash *HashInstance) Get(key object.EmeraldValue) object.EmeraldValue {
	return hash.Values.GetOrDefault(key.HashKey(), NULL)
}

func (hash *HashInstance) Set(key object.EmeraldValue, value object.EmeraldValue) {
	hashKey := key.HashKey()
	hash.Values.Set(hashKey, value)
	hash.Keys[hashKey] = key
}

func (hash *HashInstance) Each(callback func(key object.EmeraldValue, value object.EmeraldValue)) {
	for el := hash.Values.Front(); el != nil; el = el.Next() {
		callback(hash.Keys[el.Key], el.Value)
	}
}

//func hashToS() object.BuiltInMethod {
//	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
//		if _, err := EnforceArity(args, kwargs, 0, 0, []string{}); err != nil {
//			return err
//		}
//
//		var out strings.Builder
//
//		out.WriteRune('{')
//
//		ctx.Self.(*HashInstance).Each(func(key object.EmeraldValue, value object.EmeraldValue) {
//			out.WriteString(Send(key, "to_s", NULL).(*StringInstance).Value)
//			out.WriteString(Send(value, "to_s", NULL).(*StringInstance).Value)
//		})
//
//		out.WriteRune('}')
//
//		return NewString(out.String())
//	}
//}

func hashIndexAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 1, 1, []string{}); err != nil {
			return err
		}

		return ctx.Self.(*HashInstance).Get(args[0])
	}
}

func hashIndexSetter() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 2, 2, []string{}); err != nil {
			return err
		}

		ctx.Self.(*HashInstance).Set(args[0], args[1])

		return args[1]
	}
}

func hashEach() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		hash := ctx.Self.(*HashInstance)

		hash.Each(func(key object.EmeraldValue, value object.EmeraldValue) {
			ctx.Yield(key, value)
		})

		return hash
	}
}

func hashEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		args, err := EnforceArity(args, kwargs, 1, 1, []string{})

		if err != nil {
			return err
		}

		otherObj := args[0]
		if otherObj.Class().Super() != Hash {
			return FALSE
		}

		hash := ctx.Self.(*HashInstance)
		otherHash := otherObj.(*HashInstance)

		for el := hash.Values.Front(); el != nil; el = el.Next() {
			otherValue, ok := otherHash.Values.Get(el.Key)
			if !ok {
				return FALSE
			}

			if Send(el.Value, "==", NULL, otherValue) != TRUE {
				return FALSE
			}
		}

		return TRUE
	}
}
