package core

import (
	"emerald/object"
	"fmt"
	"math"
	"strconv"
)

func (rt *Runtime) NewInteger(val int64) object.EmeraldValue {
	return object.EmeraldValue{TypeID: object.INTEGER_VALUE, Heap: rt.Integer.Heap, Num: uint64(val)}
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
		val := int64(ctx.Self.Num)

		return rt.NewString(strconv.FormatInt(val, 10))
	}
}

func (rt *Runtime) integerCaseEq() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return object.NewHeapObject(err)
		}

		self := int64(ctx.Self.Num)

		if args[0].Is(object.INTEGER_VALUE) {
			return rt.NativeBoolToBooleanObject(int64(args[0].Num) == self)
		} else if args[0].Is(object.FLOAT_VALUE) {
			return rt.NativeBoolToBooleanObject(int64(math.Float64frombits(args[0].Num)) == self)
		} else {
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
		return rt.NewInteger(-int64(ctx.Self.Num))
	}
}

func (rt *Runtime) integerTimes() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for i := int64(0); i < int64(ctx.Self.Num); i++ {
			ctx.Yield(map[string]object.EmeraldValue{}, rt.NewInteger(i))
			if rt.ExceptionIsRaised() {
				return rt.NULL
			}
		}

		return ctx.Self
	}
}

func (rt *Runtime) integerToF() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewFloat(float64(int64(ctx.Self.Num)))
	}
}

func (rt *Runtime) integerSpaceship() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		left := int64(ctx.Self.Num)

		if args[0].Is(object.INTEGER_VALUE) {
			right := int64(args[0].Num)
			var result int64

			diff := left - right

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
			return object.NewHeapObject(err)
		}

		if !args[0].Is(object.INTEGER_VALUE) {
			err := rt.NewTypeError(fmt.Sprintf("no implicit conversion of %s into Integer", args[0].Class().Inspect()))
			return object.NewHeapObject(rt.Raise(err))
		}

		return cb(int64(ctx.Self.Num), int64(args[0].Num))
	}
}
