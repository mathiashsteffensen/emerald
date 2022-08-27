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

func InitInteger() {
	Integer = DefineClass(Object, "Integer", Object)

	Integer.Include(Comparable)

	DefineMethod(Integer, "to_s", integerToS())
	DefineMethod(Integer, "inspect", integerToS())
	DefineMethod(Integer, "<=>", integerSpaceship())
	DefineMethod(Integer, "==", integerEquals())
	DefineMethod(Integer, "!=", integerNotEquals())
	DefineMethod(Integer, "+", integerAdd())
	DefineMethod(Integer, "-", integerSubtract())
	DefineMethod(Integer, "*", integerMultiply())
	DefineMethod(Integer, "/", integerDivide())
	DefineMethod(Integer, "times", integerTimes())
}

func integerToS() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		val := self.(*IntegerInstance).Value

		return NewString(strconv.Itoa(int(val)))
	}
}

func integerAdd() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("+", args)
		if err != nil {
			return err
		}

		return NewInteger(self.(*IntegerInstance).Value + otherVal.Value)
	}
}

func integerSubtract() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("-", args)
		if err != nil {
			return err
		}

		return NewInteger(self.(*IntegerInstance).Value - otherVal.Value)
	}
}

func integerMultiply() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("*", args)
		if err != nil {
			return err
		}

		return NewInteger(self.(*IntegerInstance).Value * otherVal.Value)
	}
}

func integerDivide() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		right, err := requireOneIntegerArg("/", args)
		if err != nil {
			return err
		}

		left := self.(*IntegerInstance)

		if left.Value%right.Value == 0 {
			return NewInteger(left.Value / right.Value)
		} else {
			return NewFloat(float64(left.Value) / float64(right.Value))
		}
	}
}

func integerEquals() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("==", args)
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(self.(*IntegerInstance).Value == otherVal.Value)
	}
}

func integerNotEquals() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		otherVal, err := requireOneIntegerArg("!=", args)
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(self.(*IntegerInstance).Value != otherVal.Value)
	}
}

func integerTimes() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		for i := int64(0); i < self.(*IntegerInstance).Value; i++ {
			yield(block, NewInteger(i))
		}

		return self
	}
}

func integerSpaceship() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		left := target.(*IntegerInstance)
		if right, ok := args[0].(*IntegerInstance); ok {
			var result int64

			diff := left.Value - right.Value

			if diff < 0 {
				result = -1
			} else if diff > 0 {
				result = 1
			} else {
				result = 0
			}

			return NewInteger(result)
		} else {
			return NULL
		}
	}
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
