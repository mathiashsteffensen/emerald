package core

import (
	"emerald/object"
	"math"
)

func (rt *Runtime) InitNumeric() {
	rt.Numeric = rt.DefineClass("Numeric", rt.Object)
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
