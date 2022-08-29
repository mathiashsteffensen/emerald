package core

import (
	"emerald/heap"
	"emerald/object"
	"fmt"
	"github.com/dlclark/regexp2"
)

var MatchData *object.Class

type MatchDataInstance struct {
	*object.Instance
	Regexp *RegexpInstance
	Match  *regexp2.Match
	Groups []regexp2.Group
}

func InitMatchData() {
	MatchData = DefineClass(Object, "MatchData", Object)

	DefineMethod(MatchData, "[]", matchDataIndexAccessor())
	DefineMethod(MatchData, "to_s", matchDataToS())
	DefineMethod(MatchData, "to_a", matchDataToA())
	DefineMethod(MatchData, "captures", matchDataCaptures())
	DefineMethod(MatchData, "regexp", matchDataRegexp())
}

func NewMatchData(regexp *RegexpInstance, match *regexp2.Match) *MatchDataInstance {
	instance := &MatchDataInstance{
		Instance: MatchData.New(),
		Regexp:   regexp,
		Match:    match,
		Groups:   match.Groups(),
	}

	heap.SetGlobalVariableString("$~", instance)
	heap.SetGlobalVariableString("$&", NewString(instance.Groups[0].String()))

	for i, group := range instance.Groups[1:] {
		heap.SetGlobalVariableString(fmt.Sprintf("$%d", i+1), NewString(group.String()))
	}

	return instance
}

func matchDataToS() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(self.(*MatchDataInstance).Groups[0].String())
	}
}

func matchDataToA() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		groups := self.(*MatchDataInstance).Groups
		captures := []object.EmeraldValue{}

		for _, group := range groups {
			captures = append(captures, NewString(group.String()))
		}

		return NewArray(captures)
	}
}

func matchDataIndexAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		groups := self.(*MatchDataInstance).Groups
		index := args[0].(*IntegerInstance).Value

		if index > int64(len(groups)-1) {
			return NULL
		}

		return NewString(groups[index].String())
	}
}

func matchDataCaptures() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		groups := self.(*MatchDataInstance).Groups[1:]
		captures := []object.EmeraldValue{}

		for _, group := range groups {
			captures = append(captures, NewString(group.String()))
		}

		return NewArray(captures)
	}
}

func matchDataRegexp() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return self.(*MatchDataInstance).Regexp
	}
}
