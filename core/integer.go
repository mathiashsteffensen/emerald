package core

import (
	"emerald/object"
	"strconv"
)

type IntegerInstance struct {
	*object.Instance
	Value int64
}

func (i *IntegerInstance) Inspect() string {
	return strconv.Itoa(int(i.Value))
}

func (rt *Runtime) NewInteger(val int64) object.EmeraldValue {
	return &IntegerInstance{rt.Integer.New(), val}
}

func (rt *Runtime) InitInteger() {
	rt.Integer = rt.DefineClass("Integer", rt.Numeric)

	rt.Integer.Include(rt.Comparable)

	rt.DefineMethod(rt.Integer, "to_s", rt.integerToS())
	rt.DefineMethod(rt.Integer, "inspect", rt.integerToS())
	rt.DefineMethod(rt.Integer, "<=>", rt.integerSpaceship())
	rt.DefineMethod(rt.Integer, "===", rt.integerCaseEq())
	rt.DefineMethod(rt.Integer, "==", rt.integerEquals())
	rt.DefineMethod(rt.Integer, "!=", rt.integerNotEquals())
	rt.DefineMethod(rt.Integer, "+", rt.integerAdd())
	rt.DefineMethod(rt.Integer, "-", rt.integerSubtract())
	rt.DefineMethod(rt.Integer, "*", rt.integerMultiply())
	rt.DefineMethod(rt.Integer, "/", rt.integerDivide())
	rt.DefineMethod(rt.Integer, "-@", rt.integerNegate())
	rt.DefineMethod(rt.Integer, "to_f", rt.integerToF())
	rt.DefineMethod(rt.Integer, "times", rt.integerTimes())
}

func (rt *Runtime) integerToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		val := ctx.Self.(*IntegerInstance).Value

		return rt.NewString(strconv.Itoa(int(val)))
	}
}

func (rt *Runtime) integerCaseEq() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		self := ctx.Self.(*IntegerInstance)

		switch other := args[0].(type) {
		case *IntegerInstance:
			return rt.NativeBoolToBooleanObject(other.Value == self.Value)
		case *FloatInstance:
			return rt.NativeBoolToBooleanObject(int64(other.Value) == self.Value)
		default:
			return rt.FALSE
		}
	}
}

func (rt *Runtime) integerAdd() object.BuiltInMethod {
	return rt.integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
		return rt.NewInteger(left + right)
	})
}

func (rt *Runtime) integerSubtract() object.BuiltInMethod {
	return rt.integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
		return rt.NewInteger(left - right)
	})
}

func (rt *Runtime) integerMultiply() object.BuiltInMethod {
	return rt.integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
		return rt.NewInteger(left * right)
	})
}

func (rt *Runtime) integerDivide() object.BuiltInMethod {
	return rt.integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
		if left%right == 0 {
			return rt.NewInteger(left / right)
		} else {
			return rt.NewFloat(float64(left) / float64(right))
		}
	})
}

func (rt *Runtime) integerEquals() object.BuiltInMethod {
	return rt.integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
		return rt.NativeBoolToBooleanObject(left == right)
	})
}

func (rt *Runtime) integerNotEquals() object.BuiltInMethod {
	return rt.integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
		return rt.NativeBoolToBooleanObject(left != right)
	})
}

func (rt *Runtime) integerNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewInteger(-ctx.Self.(*IntegerInstance).Value)
	}
}

func (rt *Runtime) integerTimes() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for i := int64(0); i < ctx.Self.(*IntegerInstance).Value; i++ {
			ctx.Yield(map[string]object.EmeraldValue{}, rt.NewInteger(i))
		}

		return ctx.Self
	}
}

func (rt *Runtime) integerToF() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewFloat(float64(ctx.Self.(*IntegerInstance).Value))
	}
}

func (rt *Runtime) integerSpaceship() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := ctx.Self.(*IntegerInstance)

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

			return rt.NewInteger(result)
		} else {
			return rt.NULL
		}
	}
}

func (rt *Runtime) integerInfixOperator(cb func(left int64, right int64) object.EmeraldValue) object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}
		right, err := EnforceArgumentType[*IntegerInstance](rt, rt.Integer, args[0])
		if err != nil {
			return err
		}

		return cb(ctx.Self.(*IntegerInstance).Value, right.Value)
	}
}
