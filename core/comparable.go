package core

import "emerald/object"

// CRuby docs for Comparable module https://ruby-doc.org/core-3.1.2/Comparable.html

var Comparable *object.Module

func InitComparable() {
	Comparable = object.NewModule(
		"Comparable",
		object.BuiltInMethodSet{
			"==": comparableEquals(),
			"<":  comparableLessThan(),
			">":  comparableGreaterThan(),
			"<=": comparableLessThanOrEquals(),
			">=": comparableGreaterThanOrEquals(),
		},
		object.BuiltInMethodSet{},
	)
}

func comparableEquals() object.BuiltInMethod {
	return comparableMethod(func(i int64) object.EmeraldValue {
		return NativeBoolToBooleanObject(i == 0)
	})
}

func comparableLessThan() object.BuiltInMethod {
	return comparableMethod(func(i int64) object.EmeraldValue {
		return NativeBoolToBooleanObject(i < 0)
	})
}

func comparableLessThanOrEquals() object.BuiltInMethod {
	return comparableMethod(func(i int64) object.EmeraldValue {
		return NativeBoolToBooleanObject(i <= 0)
	})
}

func comparableGreaterThanOrEquals() object.BuiltInMethod {
	return comparableMethod(func(i int64) object.EmeraldValue {
		return NativeBoolToBooleanObject(i >= 0)
	})
}

func comparableGreaterThan() object.BuiltInMethod {
	return comparableMethod(func(i int64) object.EmeraldValue {
		return NativeBoolToBooleanObject(i > 0)
	})
}

func comparableMethod(spaceshipCallback func(int64) object.EmeraldValue) object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		spaceshipResult, err := target.SEND(ctx, yield, "<=>", target, NULL, args...)
		if err != nil {
			return NewStandardError(err.Error())
		}

		if IsNull(spaceshipResult) {
			return spaceshipResult
		}

		return spaceshipCallback(spaceshipResult.(*IntegerInstance).Value)
	}
}
