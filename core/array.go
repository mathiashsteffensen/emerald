package core

import (
	"bytes"
	"emerald/object"
)

var Array *object.Class

type ArrayInstance struct {
	*object.Instance
	Value []object.EmeraldValue
}

func InitArray() {
	Array = DefineClass(Object, "Array", Object)

	Array.Include(Enumerable)

	DefineMethod(Array, "[]", arrayIndexAccessor())
	DefineMethod(Array, "==", arrayEquals())
	DefineMethod(Array, "push", arrayPush())
	DefineMethod(Array, "pop", arrayPop())
	DefineMethod(Array, "each", arrayEach())
	DefineMethod(Array, "to_s", arrayToS())
	DefineMethod(Array, "inspect", arrayToS())
}

func NewArray(val []object.EmeraldValue) *ArrayInstance {
	return &ArrayInstance{
		Instance: Array.New(),
		Value:    val,
	}
}

func arrayIndexAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, 1, 1); err != nil {
			return err
		}
		intArg, err := EnforceArgumentType[*IntegerInstance](Integer, args[0])

		if err != nil {
			return err
		}

		arr := ctx.Self.(*ArrayInstance).Value

		index := intArg.Value

		if index >= int64(len(arr)) {
			return NULL
		}

		return arr[index]
	}
}

func arrayPush() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.(*ArrayInstance)

		arr.Value = append(arr.Value, args...)

		return arr
	}
}

func arrayPop() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.(*ArrayInstance)

		if len(arr.Value) == 0 {
			return NULL
		}

		index := len(arr.Value) - 1
		element := arr.Value[index]

		arr.Value = arr.Value[:index]

		return element
	}
}

func arrayEach() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.(*ArrayInstance)

		for _, val := range arr.Value {
			ctx.Yield(val)
		}

		return arr
	}
}

func arrayEquals() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		arr := ctx.Self.(*ArrayInstance)
		otherArr, ok := args[0].(*ArrayInstance)
		if !ok {
			return FALSE
		}

		if len(arr.Value) != len(otherArr.Value) {
			return FALSE
		}

		for i, value := range arr.Value {
			if !IsTruthy(Send(value, "==", NULL, otherArr.Value[i])) {
				return FALSE
			}
		}

		return TRUE
	}
}

func arrayToS() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		var out bytes.Buffer

		out.WriteString("[")

		values := ctx.Self.(*ArrayInstance).Value
		for i, value := range values {
			out.WriteString(Send(value, "inspect", NULL).Inspect())

			if i != len(values)-1 {
				out.WriteString(", ")
			}
		}

		out.WriteString("]")

		return NewString(out.String())
	}
}
