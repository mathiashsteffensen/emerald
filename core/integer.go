package core

import (
	"emerald/object"
	"fmt"
	"strconv"
)

var Integer *object.Class

type IntegerInstance struct {
	*object.Instance
	Value int64
}

func (i *IntegerInstance) Inspect() string {
	return strconv.Itoa(int(i.Value))
}

func NewInteger(val int64) object.EmeraldValue {
	return &IntegerInstance{Integer.New(), val}
}

func init() {
	Integer = object.NewClass("Integer", Object, integerBuiltInMethodSet, object.BuiltInMethodSet{})
}

var integerBuiltInMethodSet = object.BuiltInMethodSet{
	"+": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("+", args)
		if err != nil {
			return err
		}

		return NewInteger(target.(*IntegerInstance).Value + otherVal.Value)
	},
	"-": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("-", args)
		if err != nil {
			return err
		}

		return NewInteger(target.(*IntegerInstance).Value - otherVal.Value)
	},
	"*": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("*", args)
		if err != nil {
			return err
		}

		return NewInteger(target.(*IntegerInstance).Value * otherVal.Value)
	},
	"/": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("/", args)
		if err != nil {
			return err
		}

		return NewInteger(target.(*IntegerInstance).Value / otherVal.Value)
	},
	"<": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("<", args)
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(target.(*IntegerInstance).Value < otherVal.Value)
	},
	">": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg(">", args)
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(target.(*IntegerInstance).Value > otherVal.Value)
	},
	"==": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("==", args)
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(target.(*IntegerInstance).Value == otherVal.Value)
	},
	"!=": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("!=", args)
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(target.(*IntegerInstance).Value != otherVal.Value)
	},
	"<=": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("!=", args)
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(target.(*IntegerInstance).Value <= otherVal.Value)
	},
	">=": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("!=", args)
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(target.(*IntegerInstance).Value >= otherVal.Value)
	},
	"to_s": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		val := target.(*IntegerInstance).Value

		return NewString(strconv.Itoa(int(val)))
	},
	"times": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		for i := int64(0); i < target.(*IntegerInstance).Value; i++ {
			yield(block, NewInteger(i))
		}

		return target
	},
}

func requireOneIntegerArg(method string, args []object.EmeraldValue) (*IntegerInstance, object.EmeraldValue /* StandardError or nil */) {
	if len(args) != 1 {
		return nil, NewStandardError(fmt.Sprintf("Integer#%s expects single argument, got %d", method, len(args)))
	}

	otherVal, ok := args[0].(*IntegerInstance)

	if !ok {
		return nil, NewStandardError(fmt.Sprintf("Integer#%s can only be passed an integer, got=%s", method, args[0].Inspect()))
	}

	return otherVal, nil
}
