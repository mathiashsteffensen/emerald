package core

import (
	"emerald/object"
	"strings"
)

type ArrayInstance struct {
	*object.Instance
	Value []object.EmeraldValue
}

func (a *ArrayInstance) Remove(index int) {
	a.Value = append(a.Value[:index], a.Value[index+1:]...)
}

func (rt *Runtime) InitArray() {
	rt.Array = rt.DefineClass("Array", rt.Object)

	rt.Array.Include(rt.Enumerable)

	rt.DefineMethod(rt.Array, "[]", rt.arrayIndexAccessor())
	rt.DefineMethod(rt.Array, "==", rt.arrayEquals())
	rt.DefineMethod(rt.Array, "<<", rt.arrayPush())
	rt.DefineMethod(rt.Array, "push", rt.arrayPush())
	rt.DefineMethod(rt.Array, "pop", rt.arrayPop())
	rt.DefineMethod(rt.Array, "each", rt.arrayEach())
	rt.DefineMethod(rt.Array, "compact!", rt.arrayCompactBang())
	rt.DefineMethod(rt.Array, "to_s", rt.arrayToS())
	rt.DefineMethod(rt.Array, "inspect", rt.arrayToS())
	rt.DefineMethod(rt.Array, "length", rt.arrayLength())
	rt.DefineMethod(rt.Array, "size", rt.arrayLength())
}

func (rt *Runtime) NewArray(val []object.EmeraldValue) object.EmeraldValue {
	return object.NewHeapObject(&ArrayInstance{
		Instance: rt.Array.Heap.(*object.Class).New(),
		Value:    val,
	})
}

func (rt *Runtime) arrayIndexAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return object.NewHeapObject(err)
		}
		index, err := rt.EnforceIntegerArg(args[0])

		if err != nil {
			return object.NewHeapObject(err)
		}

		arr := ctx.Self.Heap.(*ArrayInstance).Value

		if index < 0 {
			index = int64(len(arr)) + index
		}

		if index < 0 || index >= int64(len(arr)) {
			return rt.NULL
		}

		return arr[index]
	}
}

func (rt *Runtime) arrayPush() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.Heap.(*ArrayInstance)

		arr.Value = append(arr.Value, args...)

		return ctx.Self
	}
}

func (rt *Runtime) arrayPop() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.Heap.(*ArrayInstance)

		if len(arr.Value) == 0 {
			return rt.NULL
		}

		index := len(arr.Value) - 1
		element := arr.Value[index]

		arr.Value = arr.Value[:index]

		return element
	}
}

func (rt *Runtime) arrayEach() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.Heap.(*ArrayInstance)

		for _, val := range arr.Value {
			if ctx.ExecutionError() != nil {
				return rt.NULL
			}

			ctx.Yield(map[string]object.EmeraldValue{}, val)
			if ctx.ExecutionError() != nil || rt.ExceptionIsRaised() {
				return rt.NULL
			}
		}

		return ctx.Self
	}
}

func (rt *Runtime) arrayCompactBang() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.Heap.(*ArrayInstance)

		i := 0 // output index
		for _, x := range arr.Value {
			if !x.IsNil() {
				// copy and increment index
				arr.Value[i] = x
				i++
			}
		}
		// Prevent memory leak by erasing truncated values
		for j := i; j < len(arr.Value); j++ {
			arr.Value[j] = object.EmeraldValue{}
		}
		arr.Value = arr.Value[:i]

		return ctx.Self
	}
}

func (rt *Runtime) arrayEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.Heap.(*ArrayInstance)
		otherArr, ok := args[0].Heap.(*ArrayInstance)
		if !ok {
			return rt.FALSE
		}

		if len(arr.Value) != len(otherArr.Value) {
			return rt.FALSE
		}

		for i, value := range arr.Value {
			if !rt.IsTruthy(rt.Send(value, "==", rt.NULL, map[string]object.EmeraldValue{}, otherArr.Value[i])) {
				return rt.FALSE
			}
		}

		return rt.TRUE
	}
}

func (rt *Runtime) arrayToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var out strings.Builder

		out.WriteString("[")

		values := ctx.Self.Heap.(*ArrayInstance).Value
		for i, value := range values {
			out.WriteString(rt.Send(value, "inspect", rt.NULL, map[string]object.EmeraldValue{}).Inspect())

			if i != len(values)-1 {
				out.WriteString(", ")
			}
		}

		out.WriteString("]")

		return rt.NewString(out.String())
	}
}
func (rt *Runtime) arrayLength() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.Heap.(*ArrayInstance)

		return rt.NewInteger(int64(len(arr.Value)))
	}
}
