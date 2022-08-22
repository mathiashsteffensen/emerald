package core

import (
	"emerald/object"
	"strconv"
)

var Float *object.Class

func InitFloat() {
	Float = object.NewClass("Float", Object, Object.Class(), object.BuiltInMethodSet{
		"+": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			left := target.(*FloatInstance)

			var newValue float64

			switch right := args[0].(type) {
			case *IntegerInstance:
				newValue = left.Value + float64(right.Value)
			case *FloatInstance:
				newValue = left.Value + right.Value
			}

			return NewFloat(newValue)
		},
		"-": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			left := target.(*FloatInstance)

			var newValue float64

			switch right := args[0].(type) {
			case *IntegerInstance:
				newValue = left.Value - float64(right.Value)
			case *FloatInstance:
				newValue = left.Value - right.Value
			}

			return NewFloat(newValue)
		},
		"*": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			left := target.(*FloatInstance)

			var newValue float64

			switch right := args[0].(type) {
			case *IntegerInstance:
				newValue = left.Value * float64(right.Value)
			case *FloatInstance:
				newValue = left.Value * right.Value
			}

			return NewFloat(newValue)
		},
		"/": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			left := target.(*FloatInstance)

			var newValue float64

			switch right := args[0].(type) {
			case *IntegerInstance:
				newValue = left.Value / float64(right.Value)
			case *FloatInstance:
				newValue = left.Value / right.Value
			}

			return NewFloat(newValue)
		},
		"<=>": floatSpaceship(),
		"to_s": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NewString(strconv.FormatFloat(target.(*FloatInstance).Value, 'f', -1, 64))
		},
	}, object.BuiltInMethodSet{}, Comparable)
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

func floatSpaceship() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		left := target.(*FloatInstance)
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
