package core

import (
	"emerald/object"
	"fmt"
	"time"

	"github.com/dlclark/regexp2"
)

type RegexpInstance struct {
	*object.Instance
	Source     string
	Expression *regexp2.Regexp
}

func (rt *Runtime) NewRegexp(pattern string) object.EmeraldValue {
	return object.NewHeapObject(&RegexpInstance{
		Instance:   rt.Regexp.Heap.(*object.Class).New(),
		Source:     pattern,
		Expression: regexp2.MustCompile(pattern, 0),
	})
}

func (rt *Runtime) InitRegexp() {
	rt.Regexp = rt.DefineClass("Regexp", rt.Object)

	rt.DefineSingletonMethod(rt.Regexp, "new", rt.regexpNew())
	rt.DefineSingletonMethod(rt.Regexp, "last_match", rt.regexpLastMatch())

	rt.DefineMethod(rt.Regexp, "inspect", rt.regexpInspect())
	rt.DefineMethod(rt.Regexp, "match", rt.regexpMatch())
	rt.DefineMethod(rt.Regexp, "=~", rt.regexpMatch())
	rt.DefineMethod(rt.Regexp, "===", rt.regexpCaseEqual())
}

func (rt *Runtime) regexpNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewRegexp(args[0].Heap.(*StringInstance).Value)
	}
}

func (rt *Runtime) regexpLastMatch() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		lastMatch := rt.Heap.GetGlobalVariableString("$~")
		if lastMatch.IsNil() {
			return rt.NULL
		}
		return lastMatch
	}
}

func (rt *Runtime) regexpInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(fmt.Sprintf("/%s/", ctx.Self.Heap.(*RegexpInstance).Source))
	}
}

func (rt *Runtime) regexpCaseEqual() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return object.NewHeapObject(err)
		}
		str, ok := args[0].Heap.(*StringInstance)
		if !ok {
			return rt.FALSE
		}
		return rt.NativeBoolToBooleanObject(rt.IsTruthy(rt.regexStringMatch(ctx, ctx.Self.Heap.(*RegexpInstance), str)))
	}
}

func (rt *Runtime) regexpMatch() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.regexStringMatch(ctx, ctx.Self.Heap.(*RegexpInstance), args[0].Heap.(*StringInstance))
	}
}

func (rt *Runtime) regexStringMatch(ctx *object.Context, regex *RegexpInstance, str *StringInstance) object.EmeraldValue {
	regex.Expression.MatchTimeout = regexp2.DefaultMatchTimeout
	if ctx.ExecutionContext != nil {
		if deadline, ok := ctx.ExecutionContext.Deadline(); ok {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return rt.NULL
			}
			regex.Expression.MatchTimeout = remaining
		}
	}

	if m, err := regex.Expression.FindStringMatch(str.Value); err != nil {
		if ctx.ExecutionError() != nil {
			return rt.NULL
		}
		panic(err)
	} else if m == nil {
		return rt.NULL
	} else {
		result := rt.NewMatchData(object.NewHeapObject(regex), m)

		return result
	}
}
