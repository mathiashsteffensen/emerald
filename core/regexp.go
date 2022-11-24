package core

import (
	"emerald/heap"
	"emerald/object"
	"fmt"
	"github.com/dlclark/regexp2"
)

var Regexp *object.Class

type RegexpInstance struct {
	*object.Instance
	Source     string
	Expression *regexp2.Regexp
}

func NewRegexp(pattern string) *RegexpInstance {
	return &RegexpInstance{
		Instance:   Regexp.New(),
		Source:     pattern,
		Expression: regexp2.MustCompile(pattern, 0),
	}
}

func InitRegexp() {
	Regexp = DefineClass("Regexp", Object)

	DefineSingletonMethod(Regexp, "new", regexpNew())
	DefineSingletonMethod(Regexp, "last_match", regexpLastMatch())

	DefineMethod(Regexp, "inspect", regexpInspect())
	DefineMethod(Regexp, "match", regexpMatch())
	DefineMethod(Regexp, "=~", regexpMatch())
}

func regexpNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewRegexp(args[0].(*StringInstance).Value)
	}
}

func regexpLastMatch() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		lastMatch := heap.GetGlobalVariableString("$~")
		if lastMatch == nil {
			return NULL
		}
		return lastMatch
	}
}

func regexpInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(fmt.Sprintf("/%s/", ctx.Self.(*RegexpInstance).Source))
	}
}

func regexpMatch() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return regexStringMatch(ctx.Self.(*RegexpInstance), args[0].(*StringInstance))
	}
}

func regexStringMatch(regex *RegexpInstance, str *StringInstance) object.EmeraldValue {
	if m, err := regex.Expression.FindStringMatch(str.Value); err != nil {
		panic(err)
	} else if m == nil {
		return NULL
	} else {
		result := NewMatchData(regex, m)

		return result
	}
}
