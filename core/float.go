package core

import (
	"emerald/object"
	"math"
	"strconv"
)

func (rt *Runtime) InitFloat() {
	rt.Float = rt.DefineClass("Float", rt.Numeric)

	rt.Float.Include(rt.Comparable)

	rt.DefineMethod(rt.Float, "to_s", rt.floatToS())
	rt.DefineMethod(rt.Float, "inspect", rt.floatToS())
	rt.DefineMethod(rt.Float, "<=>", rt.floatSpaceship())
	rt.DefineMethod(rt.Float, "+", rt.floatAdd())
	rt.DefineMethod(rt.Float, "-", rt.floatSubtract())
	rt.DefineMethod(rt.Float, "*", rt.floatMultiply())
	rt.DefineMethod(rt.Float, "/", rt.floatDivide())
	rt.DefineMethod(rt.Float, "-@", rt.floatNegate())
	rt.DefineMethod(rt.Float, "round", rt.floatRound())
}

type FloatInstance struct {
	*object.Instance
	Value float64
}

func (rt *Runtime) NewFloat(val float64) *FloatInstance {
	return &FloatInstance{
		Instance: rt.Float.New(),
		Value:    val,
	}
}

func (rt *Runtime) floatToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(strconv.FormatFloat(ctx.Self.(*FloatInstance).Value, 'f', -1, 64))
	}
}

func (rt *Runtime) floatAdd() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)

		var newValue float64

		switch right := args[0].(type) {
		case *IntegerInstance:
			newValue = left.Value + float64(right.Value)
		case *FloatInstance:
			newValue = left.Value + right.Value
		}

		return rt.NewFloat(newValue)
	}
}

func (rt *Runtime) floatSubtract() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)

		var newValue float64

		switch right := args[0].(type) {
		case *IntegerInstance:
			newValue = left.Value - float64(right.Value)
		case *FloatInstance:
			newValue = left.Value - right.Value
		}

		return rt.NewFloat(newValue)
	}
}

func (rt *Runtime) floatMultiply() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)

		var newValue float64

		switch right := args[0].(type) {
		case *IntegerInstance:
			newValue = left.Value * float64(right.Value)
		case *FloatInstance:
			newValue = left.Value * right.Value
		}

		return rt.NewFloat(newValue)
	}
}

func (rt *Runtime) floatDivide() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)

		var newValue float64

		switch right := args[0].(type) {
		case *IntegerInstance:
			newValue = left.Value / float64(right.Value)
		case *FloatInstance:
			newValue = left.Value / right.Value
		}

		return rt.NewFloat(newValue)
	}
}

func (rt *Runtime) floatNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewFloat(-ctx.Self.(*FloatInstance).Value)
	}
}

func (rt *Runtime) floatSpaceship() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)
		if right, ok := args[0].(*FloatInstance); ok {
			var result int64

			diff := left.Value - right.Value

			if diff < 0 {
				result = -1
			} else if diff > 0 {
				result = 1
			} else {
				result = 0
			}

			return rt.NewInteger(result)
		} else {
			return rt.NULL
		}
	}
}

func (rt *Runtime) roundFloatToPrecision(val float64, precision int64) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (rt *Runtime) floatRound() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		rt.EnforceArity(args, kwargs, 0, 1)

		float := ctx.Self.(*FloatInstance)

		if len(args) == 1 {
			precision, err := EnforceArgumentType[*IntegerInstance](rt, rt.Integer, args[0])
			if err != nil {
				return err
			}

			return rt.NewFloat(rt.roundFloatToPrecision(float.Value, precision.Value))
		} else {
			return rt.NewInteger(int64(math.Round(float.Value)))
		}
	}
}
