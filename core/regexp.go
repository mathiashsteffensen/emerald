package core

import (
	"emerald/object"
	"fmt"

	"github.com/dlclark/regexp2"
)

type RegexpInstance struct {
	*object.Instance
	Source     string
	Expression *regexp2.Regexp
}

func (rt *Runtime) NewRegexp(pattern string) *RegexpInstance {
	return &RegexpInstance{
		Instance:   rt.Regexp.New(),
		Source:     pattern,
		Expression: regexp2.MustCompile(pattern, 0),
	}
}

func (rt *Runtime) InitRegexp() {
	rt.Regexp = rt.DefineClass("Regexp", rt.Object)

	rt.DefineSingletonMethod(rt.Regexp, "new", rt.regexpNew())
	rt.DefineSingletonMethod(rt.Regexp, "last_match", rt.regexpLastMatch())

	rt.DefineMethod(rt.Regexp, "inspect", rt.regexpInspect())
	rt.DefineMethod(rt.Regexp, "match", rt.regexpMatch())
	rt.DefineMethod(rt.Regexp, "=~", rt.regexpMatch())
}

func (rt *Runtime) regexpNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewRegexp(args[0].(*StringInstance).Value)
	}
}

func (rt *Runtime) regexpLastMatch() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		lastMatch := rt.Heap.GetGlobalVariableString("$~")
		if lastMatch == nil {
			return rt.NULL
		}
		return lastMatch
	}
}

func (rt *Runtime) regexpInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(fmt.Sprintf("/%s/", ctx.Self.(*RegexpInstance).Source))
	}
}

func (rt *Runtime) regexpMatch() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.regexStringMatch(ctx.Self.(*RegexpInstance), args[0].(*StringInstance))
	}
}

func (rt *Runtime) regexStringMatch(regex *RegexpInstance, str *StringInstance) object.EmeraldValue {
	if m, err := regex.Expression.FindStringMatch(str.Value); err != nil {
		panic(err)
	} else if m == nil {
		return rt.NULL
	} else {
		result := rt.NewMatchData(regex, m)

		return result
	}
}
