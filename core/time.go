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
	var timeNew object.BuiltInMethod = func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewTime(time.Now())
	}

	Time = object.NewClass("Time", Object, Object.Class(), object.BuiltInMethodSet{
		"-": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			newVal := target.(*TimeInstance).Value.Sub(args[0].(*TimeInstance).Value)

			return NewInteger(newVal.Milliseconds())
		},
	}, object.BuiltInMethodSet{
		"now": timeNew,
		"new": timeNew,
	})
}
