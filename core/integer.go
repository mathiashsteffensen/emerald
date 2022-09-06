package core

import (
	"emerald/object"
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
	DefineMethod(Integer, "==", integerEquals)
	DefineMethod(Integer, "!=", integerNotEquals)
	DefineMethod(Integer, "+", integerAdd)
	DefineMethod(Integer, "-", integerSubtract)
	DefineMethod(Integer, "*", integerMultiply)
	DefineMethod(Integer, "/", integerDivide)
	DefineMethod(Integer, "times", integerTimes())
}

func integerToS() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		val := ctx.Self.(*IntegerInstance).Value

		return NewString(strconv.Itoa(int(val)))
	}
}

var integerAdd = integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
	return NewInteger(left + right)
})

var integerSubtract = integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
	return NewInteger(left - right)
})

var integerMultiply = integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
	return NewInteger(left * right)
})

var integerDivide = integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
	if left%right == 0 {
		return NewInteger(left / right)
	} else {
		return NewFloat(float64(left) / float64(right))
	}
})

var integerEquals = integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
	return NativeBoolToBooleanObject(left == right)
})

var integerNotEquals = integerInfixOperator(func(left int64, right int64) object.EmeraldValue {
	return NativeBoolToBooleanObject(left != right)
})

func integerTimes() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		for i := int64(0); i < ctx.Self.(*IntegerInstance).Value; i++ {
			ctx.Yield(NewInteger(i))
		}

		return ctx.Self
	}
}

func integerSpaceship() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
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

			return NewInteger(result)
		} else {
			return NULL
		}
	}
}

func integerInfixOperator(cb func(left int64, right int64) object.EmeraldValue) object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, 1, 1); err != nil {
			return err
		}
		if err := EnforceArgumentType(Integer, args[0]); err != nil {
			return err
		}

		return cb(ctx.Self.(*IntegerInstance).Value, args[0].(*IntegerInstance).Value)
	}
}
