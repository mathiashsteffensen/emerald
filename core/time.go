package core

import (
	"emerald/object"
	"time"
)

type TimeInstance struct {
	*object.Instance
	Value time.Time
}

func (time *TimeInstance) Inspect() string {
	return time.Value.Format("2006-01-02 15:04:05.000000 -0700")
}

func (rt *Runtime) NewTime(val time.Time) object.EmeraldValue {
	return object.NewHeapObject(&TimeInstance{
		Instance: rt.Time.Heap.(*object.Class).New(),
		Value:    val,
	})
}

func (rt *Runtime) InitTime() {
	rt.Time = rt.DefineClass("Time", rt.Object)

	rt.DefineSingletonMethod(rt.Time, "new", rt.timeNew())
	rt.DefineSingletonMethod(rt.Time, "now", rt.timeNew())

	rt.DefineMethod(rt.Time, "-", rt.timeSubtract())
	rt.DefineMethod(rt.Time, "to_f", rt.timeToF())
}

func (rt *Runtime) timeNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewTime(time.Now())
	}
}

func (rt *Runtime) timeToF() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewFloat(float64(ctx.Self.Heap.(*TimeInstance).Value.UnixNano()) / 1_000_000_000.0)
	}
}

func (rt *Runtime) timeSubtract() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		newVal := ctx.Self.Heap.(*TimeInstance).Value.Sub(args[0].Heap.(*TimeInstance).Value)

		return rt.NewFloat(newVal.Seconds())
	}
}
