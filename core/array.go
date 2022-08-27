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
	DefineMethod(Array, "push", arrayPush())
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
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := self.(*ArrayInstance).Value

		intArg, ok := args[0].(*IntegerInstance)
		if !ok {
			return NewNoConversionTypeError("Integer", args[0].Class().Super().(*object.Class).Name)
		}

		index := intArg.Value

		if index >= int64(len(arr)) {
			return NULL
		}

		return arr[index]
	}
}

func arrayPush() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := self.(*ArrayInstance)

		arr.Value = append(arr.Value, args[0])

		return arr
	}
}

func arrayEach() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := self.(*ArrayInstance)

		for _, val := range arr.Value {
			yield(block, val)
		}

		return arr
	}
}

func arrayToS() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		var out bytes.Buffer

		out.WriteString("[")

		values := self.(*ArrayInstance).Value
		for i, value := range values {
			out.WriteString(value.Inspect())

			if i != len(values)-1 {
				out.WriteString(", ")
			}
		}

		out.WriteString("]")

		return NewString(out.String())
	}
}
