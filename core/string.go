package core

import (
	"emerald/object"
	"fmt"
	"strings"
)

var String *object.Class

type StringInstance struct {
	*object.Instance
	Value string
}

func (s *StringInstance) Inspect() string { return s.Value }
func (s *StringInstance) HashKey() string { return s.Inspect() }

func NewString(val string) object.EmeraldValue {
	return &StringInstance{String.New(), val}
}

func InitString() {
	String = DefineClass(Object, "String", Object)

	DefineMethod(String, "to_s", stringToS(), false)
	DefineMethod(String, "inspect", stringInspect(), false)
	DefineMethod(String, "to_sym", stringToSym(), false)
	DefineMethod(String, "==", stringEquals(), false)
	DefineMethod(String, "+", stringAdd(), false)
	DefineMethod(String, "upcase", stringUpcase(), false)
}

func stringToS() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return target
	}
}

func stringInspect() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(fmt.Sprintf(`"%s"`, target.(*StringInstance).Value))
	}
}

func stringToSym() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewSymbol(target.Inspect())
	}
}

func stringEquals() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		left := target.(*StringInstance)
		right, ok := args[0].(*StringInstance)

		if ok {
			return NativeBoolToBooleanObject(left.Value == right.Value)
		} else {
			return NativeBoolToBooleanObject(left == right)
		}
	}
}

func stringAdd() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		targetString := target.(*StringInstance)

		argString, ok := args[0].(*StringInstance)
		if !ok {
			var typ string

			if args[0].Type() == object.CLASS_VALUE {
				typ = args[0].(*object.Class).Name
			} else {
				typ = args[0].Class().Super().(*object.Class).Name
			}

			return NewStandardError(fmt.Sprintf("no implicit conversion of %s to String", typ))
		}

		return NewString(targetString.Value + argString.Value)
	}
}

func stringUpcase() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(strings.ToUpper(target.(*StringInstance).Value))
	}
}
