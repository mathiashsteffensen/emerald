package core

import "emerald/object"

// CRuby docs for Comparable module https://ruby-doc.org/core-3.1.2/Comparable.html

var Comparable *object.Module

func InitComparable() {
	Comparable = DefineModule("Comparable")

	DefineMethod(Comparable, "==", comparableEquals())
	DefineMethod(Comparable, "<", comparableLessThan())
	DefineMethod(Comparable, ">", comparableGreaterThan())
	DefineMethod(Comparable, "<=", comparableLessThanOrEquals())
	DefineMethod(Comparable, ">=", comparableGreaterThanOrEquals())
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
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		spaceshipResult := Send(ctx.Self, "<=>", NULL, kwargs, args...)

		if spaceshipResult == NULL {
			return spaceshipResult
		}

		return spaceshipCallback(spaceshipResult.(*IntegerInstance).Value)
	}
}
