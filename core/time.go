package core

import (
	"emerald/object"
	"time"
)

var Time *object.Class

type TimeInstance struct {
	*object.Instance
	Value time.Time
}

func NewTime(val time.Time) *TimeInstance {
	return &TimeInstance{
		Instance: Time.New(),
		Value:    val,
	}
}

func init() {
	Time = DefineClass(Object, "Time", Object)

	DefineMethod(Time, "new", timeNew(), true)
	DefineMethod(Time, "now", timeNew(), true)

	DefineMethod(Time, "-", timeSubtract(), false)
	DefineMethod(Time, "to_f", timeToF(), false)
}

func timeNew() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewTime(time.Now())
	}
}

func timeToF() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewFloat(float64(target.(*TimeInstance).Value.UnixMilli()) / 1000.0)
	}
}

func timeSubtract() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		newVal := target.(*TimeInstance).Value.Sub(args[0].(*TimeInstance).Value)

		return NewInteger(newVal.Milliseconds())
	}
}
