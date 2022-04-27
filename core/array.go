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
	Array = object.NewClass("Array", Object, Object.Class(), arrayBuiltInMethodSet, object.BuiltInMethodSet{})
}

func NewArray(val []object.EmeraldValue) *ArrayInstance {
	return &ArrayInstance{
		Instance: Array.New(),
		Value:    val,
	}
}

var arrayBuiltInMethodSet = object.BuiltInMethodSet{
	"[]": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := target.(*ArrayInstance).Value

		intArg, ok := args[0].(*IntegerInstance)
		if !ok {
			return NewNoConversionTypeError("Integer", args[0].Class().Super().(*object.Class).Name)
		}

		index := intArg.Value

		if index >= int64(len(arr)) {
			return NULL
		}

		return arr[index]
	},
	"find":       arrayFind(),
	"find_index": arrayFindIndex(),
	"map":        arrayMap(),
	"each":       arrayEach(),
	"sum":        arraySum(),
	"to_s":       arrayToS(),
}

func arrayFind() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := target.(*ArrayInstance)

		for _, val := range arr.Value {
			if yield(block, val) == TRUE {
				return val
			}
		}

		return NULL
	}
}

func arrayFindIndex() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := target.(*ArrayInstance)

		for i, val := range arr.Value {
			if yield(block, val) == TRUE {
				return NewInteger(int64(i))
			}
		}

		return NULL
	}
}

func arrayMap() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := target.(*ArrayInstance)

		newArr := make([]object.EmeraldValue, len(arr.Value))

		for i, val := range arr.Value {
			newArr[i] = yield(block, val)
		}

		return NewArray(newArr)
	}
}

func arrayEach() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := target.(*ArrayInstance)

		for _, val := range arr.Value {
			yield(block, val)
		}

		return arr
	}
}

func arraySum() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		var err error
		arr := target.(*ArrayInstance)
		blockGiven := IsNull(block)

		var accumulated object.EmeraldValue
		if blockGiven {
			accumulated = yield(block, arr.Value[0])
		} else {
			accumulated = arr.Value[0]
		}

		for _, value := range arr.Value[1:] {
			if blockGiven {
				accumulated, err = accumulated.SEND(ctx, nil, "+", accumulated, nil, yield(block, value))
				if err != nil {
					return NewStandardError(err.Error())
				}
			} else {
				accumulated, err = accumulated.SEND(ctx, nil, "+", accumulated, nil, value)
				if err != nil {
					return NewStandardError(err.Error())
				}
			}
		}

		return accumulated
	}
}

func arrayToS() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		var out bytes.Buffer

		out.WriteString("[")

		values := target.(*ArrayInstance).Value
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
