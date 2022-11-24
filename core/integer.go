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
	Integer = DefineClass("Integer", Numeric)

	Integer.Include(Comparable)

	DefineMethod(Integer, "to_s", integerToS())
	DefineMethod(Integer, "inspect", integerToS())
	DefineMethod(Integer, "<=>", integerSpaceship())
	DefineMethod(Integer, "===", integerCaseEq())
	DefineMethod(Integer, "==", integerEquals)
	DefineMethod(Integer, "!=", integerNotEquals)
	DefineMethod(Integer, "+", integerAdd)
	DefineMethod(Integer, "-", integerSubtract)
	DefineMethod(Integer, "*", integerMultiply)
	DefineMethod(Integer, "/", integerDivide)
	DefineMethod(Integer, "-@", integerNegate())
	DefineMethod(Integer, "to_f", integerToF())
	DefineMethod(Integer, "times", integerTimes())
}

func integerToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		val := ctx.Self.(*IntegerInstance).Value

		return NewString(strconv.Itoa(int(val)))
	}
}

func integerCaseEq() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 1, 1, []string{}); err != nil {
			return err
		}

		self := ctx.Self.(*IntegerInstance)

		switch other := args[0].(type) {
		case *IntegerInstance:
			return NativeBoolToBooleanObject(other.Value == self.Value)
		case *FloatInstance:
			return NativeBoolToBooleanObject(int64(other.Value) == self.Value)
		default:
			return FALSE
		}
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

func integerNegate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewInteger(-ctx.Self.(*IntegerInstance).Value)
	}
}

func integerTimes() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for i := int64(0); i < ctx.Self.(*IntegerInstance).Value; i++ {
			ctx.Yield(NewInteger(i))
		}

		return ctx.Self
	}
}

func integerToF() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewFloat(float64(ctx.Self.(*IntegerInstance).Value))
	}
}

func integerSpaceship() object.BuiltInMethod {
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

			return NewInteger(result)
		} else {
			return NULL
		}
	}
}

func integerInfixOperator(cb func(left int64, right int64) object.EmeraldValue) object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 1, 1, []string{}); err != nil {
			return err
		}
		right, err := EnforceArgumentType[*IntegerInstance](Integer, args[0])
		if err != nil {
			return err
		}

		return cb(ctx.Self.(*IntegerInstance).Value, right.Value)
	}
}
