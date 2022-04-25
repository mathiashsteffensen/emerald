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
		"to_s": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			return NewString(strconv.FormatFloat(target.(*FloatInstance).Value, 'f', -1, 64))
		},
	}, object.BuiltInMethodSet{})
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
