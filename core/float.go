package core

import (
	"emerald/object"
	"fmt"
	"math"
	"strconv"
)

func (rt *Runtime) InitFloat() {
	rt.Float = rt.DefineClass("Float", rt.Numeric)

	rt.Float.Include(rt.Comparable)

	rt.DefineMethod(rt.Float, "to_s", rt.floatToS())
	rt.DefineMethod(rt.Float, "inspect", rt.floatToS())
	rt.DefineMethod(rt.Float, "<=>", rt.floatSpaceship())
	rt.DefineMethod(rt.Float, "+", rt.numericAdd())
	rt.DefineMethod(rt.Float, "-", rt.numericSubtract())
	rt.DefineMethod(rt.Float, "*", rt.numericMultiply())
	rt.DefineMethod(rt.Float, "/", rt.numericDivide())
	rt.DefineMethod(rt.Float, "-@", rt.floatNegate())
	rt.DefineMethod(rt.Float, "round", rt.floatRound())
}

func (rt *Runtime) NewFloat(val float64) object.EmeraldValue {
	return object.EmeraldValue{
		TypeID: object.FLOAT_VALUE,
		Heap:   rt.Float.Heap,
		Num:    math.Float64bits(val),
	}
}

func (rt *Runtime) floatToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(strconv.FormatFloat(math.Float64frombits(ctx.Self.Num), 'f', -1, 64))
	}
}

func (rt *Runtime) floatNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewFloat(-math.Float64frombits(ctx.Self.Num))
	}
}

func (rt *Runtime) floatSpaceship() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.numericSpaceship(ctx.Self, args[0])
	}
}

func (rt *Runtime) roundFloatToPrecision(val float64, precision int64) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (rt *Runtime) floatRound() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		rt.EnforceArity(args, kwargs, 0, 1)

		val := math.Float64frombits(ctx.Self.Num)

		if len(args) == 1 {
			if !args[0].Is(object.INTEGER_VALUE) {
				err := rt.NewTypeError(fmt.Sprintf("no implicit conversion of %s into Integer", args[0].Class().Inspect()))
				return object.NewHeapObject(rt.Raise(err))
			}
			precision := int64(args[0].Num)

			return rt.NewFloat(rt.roundFloatToPrecision(val, precision))
		} else {
			return rt.NewInteger(int64(math.Round(val)))
		}
	}
}
