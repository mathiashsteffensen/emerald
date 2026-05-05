package core

import (
	"emerald/object"
	"fmt"

	"github.com/dlclark/regexp2"
)

type MatchDataInstance struct {
	*object.Instance
	Regexp object.EmeraldValue
	Match  *regexp2.Match
	Groups []regexp2.Group
}

func (rt *Runtime) InitMatchData() {
	rt.MatchData = rt.DefineClass("MatchData", rt.Object)

	rt.DefineMethod(rt.MatchData, "[]", rt.matchDataIndexAccessor())
	rt.DefineMethod(rt.MatchData, "to_s", rt.matchDataToS())
	rt.DefineMethod(rt.MatchData, "to_a", rt.matchDataToA())
	rt.DefineMethod(rt.MatchData, "captures", rt.matchDataCaptures())
	rt.DefineMethod(rt.MatchData, "regexp", rt.matchDataRegexp())
}

func (rt *Runtime) NewMatchData(regexp object.EmeraldValue, match *regexp2.Match) object.EmeraldValue {
	instance := &MatchDataInstance{
		Instance: rt.MatchData.Heap.(*object.Class).New(),
		Regexp:   regexp,
		Match:    match,
		Groups:   match.Groups(),
	}

	val := object.NewHeapObject(instance)

	rt.Heap.SetGlobalVariableString("$~", val)
	rt.Heap.SetGlobalVariableString("$&", rt.NewString(instance.Groups[0].String()))

	for i, group := range instance.Groups[1:] {
		rt.Heap.SetGlobalVariableString(fmt.Sprintf("$%d", i+1), rt.NewString(group.String()))
	}

	return val
}

func (rt *Runtime) matchDataToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(ctx.Self.Heap.(*MatchDataInstance).Groups[0].String())
	}
}

func (rt *Runtime) matchDataToA() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		groups := ctx.Self.Heap.(*MatchDataInstance).Groups
		captures := []object.EmeraldValue{}

		for _, group := range groups {
			captures = append(captures, rt.NewString(group.String()))
		}

		return rt.NewArray(captures)
	}
}

func (rt *Runtime) matchDataIndexAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		groups := ctx.Self.Heap.(*MatchDataInstance).Groups
		index, err := rt.EnforceIntegerArg(args[0])
		if err != nil {
			return object.NewHeapObject(err)
		}

		if index > int64(len(groups)-1) {
			return rt.NULL
		}

		return rt.NewString(groups[index].String())
	}
}

func (rt *Runtime) matchDataCaptures() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		groups := ctx.Self.Heap.(*MatchDataInstance).Groups[1:]
		captures := []object.EmeraldValue{}

		for _, group := range groups {
			captures = append(captures, rt.NewString(group.String()))
		}

		return rt.NewArray(captures)
	}
}

func (rt *Runtime) matchDataRegexp() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return ctx.Self.Heap.(*MatchDataInstance).Regexp
	}
}
