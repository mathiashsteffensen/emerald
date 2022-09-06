package core

import (
	"bytes"
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

	DefineMethod(String, "to_s", stringToS())
	DefineMethod(String, "inspect", stringInspect())
	DefineMethod(String, "to_sym", stringToSym())
	DefineMethod(String, "==", stringEquals())
	DefineMethod(String, "+", stringAdd())
	DefineMethod(String, "*", stringMultiply())
	DefineMethod(String, "upcase", stringUpcase())
	DefineMethod(String, "match", stringMatch())
	DefineMethod(String, "=~", stringMatch())
}

func stringToS() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return ctx.Self
	}
}

func stringInspect() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(fmt.Sprintf(`"%s"`, ctx.Self.(*StringInstance).Value))
	}
}

func stringToSym() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewSymbol(ctx.Self.Inspect())
	}
}

func stringEquals() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*StringInstance)
		right, ok := args[0].(*StringInstance)

		if ok {
			return NativeBoolToBooleanObject(left.Value == right.Value)
		} else {
			return NativeBoolToBooleanObject(left == right)
		}
	}
}

func stringAdd() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		selfString := ctx.Self.(*StringInstance)

		if _, err := EnforceArity(args, 1, 1); err != nil {
			return NULL
		}

		if err := EnforceArgumentType(String, args[0]); err != nil {
			return err
		}

		argString := args[0].(*StringInstance)

		return NewString(selfString.Value + argString.Value)
	}
}

func stringMultiply() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		selfString := ctx.Self.(*StringInstance)

		if _, err := EnforceArity(args, 1, 1); err != nil {
			return err
		}
		if err := EnforceArgumentType(Integer, args[0]); err != nil {
			return err
		}

		arg := args[0].(*IntegerInstance)

		var newString bytes.Buffer

		for i := int64(0); i < arg.Value; i++ {
			newString.WriteString(selfString.Value)
		}

		return NewString(newString.String())
	}
}

func stringUpcase() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(strings.ToUpper(ctx.Self.(*StringInstance).Value))
	}
}

func stringMatch() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return regexStringMatch(args[0].(*RegexpInstance), ctx.Self.(*StringInstance))
	}
}
