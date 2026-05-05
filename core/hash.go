package core

import (
	"emerald/object"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

func (rt *Runtime) InitHash() {
	rt.Hash = rt.DefineClass("Hash", rt.Object)

	rt.Hash.Include(rt.Enumerable)

	rt.DefineMethod(rt.Hash, "[]", rt.hashIndexAccessor())
	rt.DefineMethod(rt.Hash, "[]=", rt.hashIndexSetter())
	rt.DefineMethod(rt.Hash, "==", rt.hashEquals())
	rt.DefineMethod(rt.Hash, "each", rt.hashEach())
	rt.DefineMethod(rt.Hash, "to_s", rt.hashToS())
	rt.DefineMethod(rt.Hash, "inspect", rt.hashToS())
}

type HashInstance struct {
	*object.Instance
	Values *orderedmap.OrderedMap[string, object.EmeraldValue] // Only Values need to be ordered since we always iterate on those
	Keys   map[string]object.EmeraldValue
}

func (rt *Runtime) NewHash() object.EmeraldValue {
	return object.NewHeapObject(&HashInstance{
		Instance: rt.Hash.Heap.(*object.Class).New(),
		Values:   orderedmap.NewOrderedMap[string, object.EmeraldValue](),
		Keys:     map[string]object.EmeraldValue{},
	})
}

func (hash *HashInstance) Get(key object.EmeraldValue) object.EmeraldValue {
	val, ok := hash.Values.Get(key.HashKey())
	if !ok {
		return object.EmeraldValue{}
	}
	return val
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

func (rt *Runtime) hashToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 0, 0); err != nil {
			return object.NewHeapObject(err)
		}

		pairs := []string{}

		ctx.Self.Heap.(*HashInstance).Each(func(key object.EmeraldValue, value object.EmeraldValue) {
			var out strings.Builder

			out.WriteString(rt.Send(key, "to_s", rt.NULL, map[string]object.EmeraldValue{}).Inspect())
			out.WriteString(" => ")
			out.WriteString(rt.Send(value, "to_s", rt.NULL, map[string]object.EmeraldValue{}).Inspect())

			pairs = append(pairs, out.String())
		})

		var out strings.Builder

		out.WriteRune('{')

		out.WriteString(strings.Join(pairs, ", "))

		out.WriteRune('}')

		return rt.NewString(out.String())
	}
}

func (rt *Runtime) hashIndexAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return object.NewHeapObject(err)
		}

		val := ctx.Self.Heap.(*HashInstance).Get(args[0])
		if val.IsNil() {
			return rt.NULL
		}
		return val
	}
}

func (rt *Runtime) hashIndexSetter() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 2, 2); err != nil {
			return object.NewHeapObject(err)
		}

		ctx.Self.Heap.(*HashInstance).Set(args[0], args[1])

		return args[1]
	}
}

func (rt *Runtime) hashEach() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		hash := ctx.Self.Heap.(*HashInstance)

		hash.Each(func(key object.EmeraldValue, value object.EmeraldValue) {
			ctx.Yield(map[string]object.EmeraldValue{}, key, value)
		})

		return ctx.Self
	}
}

func (rt *Runtime) hashEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		args, err := rt.EnforceArity(args, kwargs, 1, 1)

		if err != nil {
			return object.NewHeapObject(err)
		}

		otherObj := args[0]
		if object.RealClass(otherObj) != rt.Hash {
			return rt.FALSE
		}

		hash := ctx.Self.Heap.(*HashInstance)
		otherHash := otherObj.Heap.(*HashInstance)

		for el := hash.Values.Front(); el != nil; el = el.Next() {
			otherValue, ok := otherHash.Values.Get(el.Key)
			if !ok {
				return rt.FALSE
			}

			if rt.Send(el.Value, "==", rt.NULL, map[string]object.EmeraldValue{}, otherValue) != rt.TRUE {
				return rt.FALSE
			}
		}

		return rt.TRUE
	}
}
