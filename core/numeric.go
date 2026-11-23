package core

import (
	"emerald/object"
	"fmt"
	"math"
)

func (rt *Runtime) InitNumeric() {
	rt.Numeric = rt.DefineClass("Numeric", rt.Object)
}

func (rt *Runtime) numericAdd() object.BuiltInMethod {
	return rt.numericArithmetic(
		func(left int64, right int64) object.EmeraldValue { return rt.NewInteger(left + right) },
		func(left float64, right float64) float64 { return left + right },
	)
}

func (rt *Runtime) numericSubtract() object.BuiltInMethod {
	return rt.numericArithmetic(
		func(left int64, right int64) object.EmeraldValue { return rt.NewInteger(left - right) },
		func(left float64, right float64) float64 { return left - right },
	)
}

func (rt *Runtime) numericMultiply() object.BuiltInMethod {
	return rt.numericArithmetic(
		func(left int64, right int64) object.EmeraldValue { return rt.NewInteger(left * right) },
		func(left float64, right float64) float64 { return left * right },
	)
}

func (rt *Runtime) numericDivide() object.BuiltInMethod {
	return rt.numericArithmetic(
		func(left int64, right int64) object.EmeraldValue {
			if left%right == 0 {
				return rt.NewInteger(left / right)
			} else {
				return rt.NewFloat(float64(left) / float64(right))
			}
		},
		func(left float64, right float64) float64 { return left / right },
	)
}

// numericArithmetic builds an operator that Integer and Float share. As in Ruby,
// two Integers use integer arithmetic and a Float on either side gives a Float.
func (rt *Runtime) numericArithmetic(
	integerOperation func(left int64, right int64) object.EmeraldValue,
	floatOperation func(left float64, right float64) float64,
) object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return object.NewHeapObject(err)
		}

		left, right := ctx.Self, args[0]

		if left.Is(object.INTEGER_VALUE) && right.Is(object.INTEGER_VALUE) {
			return integerOperation(int64(left.Num), int64(right.Num))
		}

		if !right.Is(object.INTEGER_VALUE) && !right.Is(object.FLOAT_VALUE) {
			return object.NewHeapObject(rt.Raise(rt.coercionFailedError(left, right)))
		}

		return rt.NewFloat(floatOperation(numberToFloat64(left), numberToFloat64(right)))
	}
}

func numberToFloat64(number object.EmeraldValue) float64 {
	if number.Is(object.INTEGER_VALUE) {
		return float64(int64(number.Num))
	}

	return math.Float64frombits(number.Num)
}

func (rt *Runtime) coercionFailedError(self object.EmeraldValue, other object.EmeraldValue) object.EmeraldError {
	otherDescription := other.Class().Inspect()
	if other.IsImmediate() {
		otherDescription = other.Inspect()
	}

	return rt.NewTypeError(fmt.Sprintf("%s can't be coerced into %s", otherDescription, self.Class().Inspect()))
}

func (rt *Runtime) numericSpaceship(left object.EmeraldValue, right object.EmeraldValue) object.EmeraldValue {
	switch {
	case left.Is(object.INTEGER_VALUE) && right.Is(object.INTEGER_VALUE):
		return rt.NewInteger(compareInt64(int64(left.Num), int64(right.Num)))
	case left.Is(object.INTEGER_VALUE) && right.Is(object.FLOAT_VALUE):
		result, ok := compareFloat64(float64(int64(left.Num)), math.Float64frombits(right.Num))
		if !ok {
			return rt.NULL
		}
		return rt.NewInteger(result)
	case left.Is(object.FLOAT_VALUE) && right.Is(object.INTEGER_VALUE):
		result, ok := compareFloat64(math.Float64frombits(left.Num), float64(int64(right.Num)))
		if !ok {
			return rt.NULL
		}
		return rt.NewInteger(result)
	case left.Is(object.FLOAT_VALUE) && right.Is(object.FLOAT_VALUE):
		result, ok := compareFloat64(math.Float64frombits(left.Num), math.Float64frombits(right.Num))
		if !ok {
			return rt.NULL
		}
		return rt.NewInteger(result)
	default:
		return rt.NULL
	}
}

func compareInt64(left int64, right int64) int64 {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func compareFloat64(left float64, right float64) (int64, bool) {
	if math.IsNaN(left) || math.IsNaN(right) {
		return 0, false
	}

	switch {
	case left < right:
		return -1, true
	case left > right:
		return 1, true
	default:
		return 0, true
	}
}
