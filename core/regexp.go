package core

import (
	"emerald/object"
	"github.com/dlclark/regexp2"
)

var Regexp *object.Class

type RegexpInstance struct {
	*object.Instance
	expression *regexp2.Regexp
}

func NewRegexp(pattern string) *RegexpInstance {
	return &RegexpInstance{
		Instance:   Regexp.New(),
		expression: regexp2.MustCompile(pattern, 0),
	}
}

func InitRegexp() {
	Regexp = DefineClass(Object, "Regexp", Object)

	DefineSingletonMethod(Regexp, "new", regexpNew())

	DefineMethod(Regexp, "match", regexpMatch())
	DefineMethod(Regexp, "=~", regexpMatch())
}

func regexpNew() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewRegexp(args[0].(*StringInstance).Value)
	}
}

func regexpMatch() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return regexStringMatch(self.(*RegexpInstance), args[0].(*StringInstance))
	}
}

func regexStringMatch(regex *RegexpInstance, str *StringInstance) object.EmeraldValue {
	if isMatch, _ := regex.expression.MatchString(str.Value); isMatch {
		return TRUE
	}

	return FALSE
}
