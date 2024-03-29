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
	String = DefineClass("String", Object)

	DefineSingletonMethod(String, "new", stringNew())

	DefineMethod(String, "to_s", stringToS())
	DefineMethod(String, "inspect", stringInspect())
	DefineMethod(String, "to_sym", stringToSym())
	DefineMethod(String, "==", stringEquals())
	DefineMethod(String, "+", stringAdd())
	DefineMethod(String, "*", stringMultiply())
	DefineMethod(String, "=~", stringMatch())
	DefineMethod(String, "match", stringMatch())
	DefineMethod(String, "upcase", stringUpcase())
	DefineMethod(String, "size", stringSize())
	DefineMethod(String, "length", stringSize())
	DefineMethod(String, "split", stringSplit())
}

func stringNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		args, err := EnforceArity(args, kwargs, 0, 1)
		if err != nil {
			return err
		}

		if len(args) == 1 {
			str, err := EnforceArgumentType[*StringInstance](String, args[0])
			if err != nil {
				return err
			}

			return str
		}

		return NewString("")
	}
}

func stringToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return ctx.Self
	}
}

func stringInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(fmt.Sprintf("%q", ctx.Self.(*StringInstance).Value))
	}
}

func stringToSym() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewSymbol(ctx.Self.Inspect())
	}
}

func stringEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
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
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		selfString := ctx.Self.(*StringInstance)

		if _, err := EnforceArity(args, kwargs, 1, 1); err != nil {
			return NULL
		}

		str, err := EnforceArgumentType[*StringInstance](String, args[0])
		if err != nil {
			return err
		}

		return NewString(selfString.Value + str.Value)
	}
}

func stringMultiply() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		selfString := ctx.Self.(*StringInstance)

		if _, err := EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		arg, err := EnforceArgumentType[*IntegerInstance](Integer, args[0])
		if err != nil {
			return err
		}

		var newString strings.Builder

		for i := int64(0); i < arg.Value; i++ {
			newString.WriteString(selfString.Value)
		}

		return NewString(newString.String())
	}
}

func stringUpcase() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(strings.ToUpper(ctx.Self.(*StringInstance).Value))
	}
}

func stringMatch() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return regexStringMatch(args[0].(*RegexpInstance), ctx.Self.(*StringInstance))
	}
}

func stringSize() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewInteger(
			int64(len(ctx.Self.(*StringInstance).Value)),
		)
	}
}

func stringSplit() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 0, 1); err != nil {
			return err
		}

		var sep string

		if len(args) == 0 {
			sep = " "
		} else {
			arg, err := EnforceArgumentType[*StringInstance](String, args[0])
			if err != nil {
				return err
			}

			sep = arg.Value
		}

		self := ctx.Self.(*StringInstance)

		slice := []object.EmeraldValue{}

		for _, s := range strings.Split(self.Value, sep) {
			slice = append(slice, NewString(s))
		}

		return NewArray(slice)
	}
}
