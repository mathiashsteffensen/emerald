package core

import (
	"emerald/object"
	"math"
	"strconv"
)

var Float *object.Class

func InitFloat() {
	Float = DefineClass("Float", Numeric)

	Float.Include(Comparable)

	DefineMethod(Float, "to_s", floatToS())
	DefineMethod(Float, "inspect", floatToS())
	DefineMethod(Float, "<=>", floatSpaceship())
	DefineMethod(Float, "+", floatAdd())
	DefineMethod(Float, "-", floatSubtract())
	DefineMethod(Float, "*", floatMultiply())
	DefineMethod(Float, "/", floatDivide())
	DefineMethod(Float, "-@", floatNegate())
	DefineMethod(Float, "round", floatRound())
}

type FloatInstance struct {
	*object.Instance
	Value float64
}

func NewFloat(val float64) *FloatInstance {
	return &FloatInstance{
		Instance: Float.New(),
		Value:    val,
	}
}

func floatToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(strconv.FormatFloat(ctx.Self.(*FloatInstance).Value, 'f', -1, 64))
	}
}

func floatAdd() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)

		var newValue float64

		switch right := args[0].(type) {
		case *IntegerInstance:
			newValue = left.Value + float64(right.Value)
		case *FloatInstance:
			newValue = left.Value + right.Value
		}

		return NewFloat(newValue)
	}
}

func floatSubtract() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)

		var newValue float64

		switch right := args[0].(type) {
		case *IntegerInstance:
			newValue = left.Value - float64(right.Value)
		case *FloatInstance:
			newValue = left.Value - right.Value
		}

		return NewFloat(newValue)
	}
}

func floatMultiply() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)

		var newValue float64

		switch right := args[0].(type) {
		case *IntegerInstance:
			newValue = left.Value * float64(right.Value)
		case *FloatInstance:
			newValue = left.Value * right.Value
		}

		return NewFloat(newValue)
	}
}

func floatDivide() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*FloatInstance)

		var newValue float64

		switch right := args[0].(type) {
		case *IntegerInstance:
			newValue = left.Value / float64(right.Value)
		case *FloatInstance:
			newValue = left.Value / right.Value
		}

		return NewFloat(newValue)
	}
}

func floatNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewFloat(-ctx.Self.(*FloatInstance).Value)
	}
}

func floatSpaceship() object.BuiltInMethod {
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

			return NewInteger(result)
		} else {
			return NULL
		}
	}
}

func roundFloatToPrecision(val float64, precision int64) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func floatRound() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		EnforceArity(args, kwargs, 0, 1)

		float := ctx.Self.(*FloatInstance)

		if len(args) == 1 {
			precision, err := EnforceArgumentType[*IntegerInstance](Integer, args[0])
			if err != nil {
				return err
			}

			return NewFloat(roundFloatToPrecision(float.Value, precision.Value))
		} else {
			return NewInteger(int64(math.Round(float.Value)))
		}
	}
}
