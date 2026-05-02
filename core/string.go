package core

import (
	"emerald/object"
	"fmt"
	"strings"
)

type StringInstance struct {
	*object.Instance
	Value string
}

func (s *StringInstance) Inspect() string { return s.Value }
func (s *StringInstance) HashKey() string { return s.Inspect() }

func (rt *Runtime) NewString(val string) object.EmeraldValue {
	return &StringInstance{rt.String.New(), val}
}

func (rt *Runtime) InitString() {
	rt.String = rt.DefineClass("String", rt.Object)

	rt.DefineSingletonMethod(rt.String, "new", rt.stringNew())

	rt.DefineMethod(rt.String, "to_s", rt.stringToS())
	rt.DefineMethod(rt.String, "inspect", rt.stringInspect())
	rt.DefineMethod(rt.String, "to_sym", rt.stringToSym())
	rt.DefineMethod(rt.String, "==", rt.stringEquals())
	rt.DefineMethod(rt.String, "+", rt.stringAdd())
	rt.DefineMethod(rt.String, "*", rt.stringMultiply())
	rt.DefineMethod(rt.String, "=~", rt.stringMatch())
	rt.DefineMethod(rt.String, "match", rt.stringMatch())
	rt.DefineMethod(rt.String, "upcase", rt.stringUpcase())
	rt.DefineMethod(rt.String, "size", rt.stringSize())
	rt.DefineMethod(rt.String, "length", rt.stringSize())
	rt.DefineMethod(rt.String, "split", rt.stringSplit())
}

func (rt *Runtime) stringNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		args, err := rt.EnforceArity(args, kwargs, 0, 1)
		if err != nil {
			return err
		}

		if len(args) == 1 {
			str, err := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
			if err != nil {
				return err
			}

			return str
		}

		return rt.NewString("")
	}
}

func (rt *Runtime) stringToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return ctx.Self
	}
}

func (rt *Runtime) stringInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(fmt.Sprintf("%q", ctx.Self.(*StringInstance).Value))
	}
}

func (rt *Runtime) stringToSym() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewSymbol(ctx.Self.Inspect())
	}
}

func (rt *Runtime) stringEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*StringInstance)
		right, ok := args[0].(*StringInstance)

		if ok {
			return rt.NativeBoolToBooleanObject(left.Value == right.Value)
		} else {
			return rt.NativeBoolToBooleanObject(left == right)
		}
	}
}

func (rt *Runtime) stringAdd() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		selfString := ctx.Self.(*StringInstance)

		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return rt.NULL
		}

		str, err := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
		if err != nil {
			return err
		}

		return rt.NewString(selfString.Value + str.Value)
	}
}

func (rt *Runtime) stringMultiply() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		selfString := ctx.Self.(*StringInstance)

		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		arg, err := EnforceArgumentType[*IntegerInstance](rt, rt.Integer, args[0])
		if err != nil {
			return err
		}

		var newString strings.Builder

		for i := int64(0); i < arg.Value; i++ {
			newString.WriteString(selfString.Value)
		}

		return rt.NewString(newString.String())
	}
}

func (rt *Runtime) stringUpcase() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(strings.ToUpper(ctx.Self.(*StringInstance).Value))
	}
}

func (rt *Runtime) stringMatch() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.regexStringMatch(args[0].(*RegexpInstance), ctx.Self.(*StringInstance))
	}
}

func (rt *Runtime) stringSize() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewInteger(
			int64(len(ctx.Self.(*StringInstance).Value)),
		)
	}
}

func (rt *Runtime) stringSplit() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 0, 1); err != nil {
			return err
		}

		var sep string

		if len(args) == 0 {
			sep = " "
		} else {
			arg, err := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
			if err != nil {
				return err
			}

			sep = arg.Value
		}

		self := ctx.Self.(*StringInstance)

		slice := []object.EmeraldValue{}

		for _, s := range strings.Split(self.Value, sep) {
			slice = append(slice, rt.NewString(s))
		}

		return rt.NewArray(slice)
	}
}
