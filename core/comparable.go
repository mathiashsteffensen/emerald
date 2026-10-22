package core

import (
	"emerald/object"
	"fmt"
)

// CRuby docs for rt.Comparable module https://ruby-doc.org/core-3.1.2/rt.Comparable.html

func (rt *Runtime) InitComparable() {
	rt.Comparable = rt.DefineModule("Comparable")

	rt.DefineMethod(rt.Comparable, "==", rt.comparableEquals())
	rt.DefineMethod(rt.Comparable, "<", rt.comparableLessThan())
	rt.DefineMethod(rt.Comparable, ">", rt.comparableGreaterThan())
	rt.DefineMethod(rt.Comparable, "<=", rt.comparableLessThanOrEquals())
	rt.DefineMethod(rt.Comparable, ">=", rt.comparableGreaterThanOrEquals())
}

func (rt *Runtime) comparableEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return object.NewHeapObject(err)
		}

		spaceshipResult := rt.Send(ctx.Self, "<=>", rt.NULL, kwargs, args...)

		return rt.NativeBoolToBooleanObject(spaceshipResult.Is(object.INTEGER_VALUE) && spaceshipResult.Num == 0)
	}
}

func (rt *Runtime) comparableLessThan() object.BuiltInMethod {
	return rt.comparableMethod(func(i int64) object.EmeraldValue {
		return rt.NativeBoolToBooleanObject(i < 0)
	})
}

func (rt *Runtime) comparableLessThanOrEquals() object.BuiltInMethod {
	return rt.comparableMethod(func(i int64) object.EmeraldValue {
		return rt.NativeBoolToBooleanObject(i <= 0)
	})
}

func (rt *Runtime) comparableGreaterThanOrEquals() object.BuiltInMethod {
	return rt.comparableMethod(func(i int64) object.EmeraldValue {
		return rt.NativeBoolToBooleanObject(i >= 0)
	})
}

func (rt *Runtime) comparableGreaterThan() object.BuiltInMethod {
	return rt.comparableMethod(func(i int64) object.EmeraldValue {
		return rt.NativeBoolToBooleanObject(i > 0)
	})
}

func (rt *Runtime) comparableMethod(spaceshipCallback func(int64) object.EmeraldValue) object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return object.NewHeapObject(err)
		}

		spaceshipResult := rt.Send(ctx.Self, "<=>", rt.NULL, kwargs, args...)

		if spaceshipResult.IsNil() {
			other := args[0].Class().Inspect()
			if args[0].IsImmediate() {
				other = args[0].Inspect()
			}

			err := rt.newArgumentError(fmt.Sprintf("comparison of %s with %s failed", ctx.Self.Class().Inspect(), other))
			return object.NewHeapObject(rt.Raise(err))
		}

		return spaceshipCallback(int64(spaceshipResult.Num))
	}
}
