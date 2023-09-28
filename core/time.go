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

func (time *TimeInstance) Inspect() string {
	return time.Value.Format("2006-01-02 15:04:05.000000 -0700")
}

func NewTime(val time.Time) *TimeInstance {
	return &TimeInstance{
		Instance: Time.New(),
		Value:    val,
	}
}

func init() {
	Time = DefineClass("Time", Object)

	DefineSingletonMethod(Time, "new", timeNew())
	DefineSingletonMethod(Time, "now", timeNew())

	DefineMethod(Time, "-", timeSubtract())
	DefineMethod(Time, "to_f", timeToF())
}

func timeNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewTime(time.Now())
	}
}

func timeToF() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewFloat(float64(ctx.Self.(*TimeInstance).Value.UnixNano()) / 1_000_000.0)
	}
}

func timeSubtract() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		newVal := ctx.Self.(*TimeInstance).Value.Sub(args[0].(*TimeInstance).Value)

		return NewInteger(newVal.Milliseconds())
	}
}
