package core

import "emerald/object"

var Array *object.Class

type ArrayInstance struct {
	*object.Instance
	Value []object.EmeraldValue
}

func init() {
	Array = object.NewClass("Array", Object, arrayBuiltInMethodSet, object.BuiltInMethodSet{})
}

func NewArray(val []object.EmeraldValue) *ArrayInstance {
	return &ArrayInstance{
		Instance: Array.New(),
		Value:    val,
	}
}

var arrayBuiltInMethodSet = object.BuiltInMethodSet{
	"find":       arrayFind(),
	"find_index": arrayFindIndex(),
	"map":        arrayMap(),
	"each":       arrayEach(),
}

func arrayFind() object.BuiltInMethod {
	return func(target object.EmeraldValue, block *object.Block, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
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
	return func(target object.EmeraldValue, block *object.Block, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
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
	return func(target object.EmeraldValue, block *object.Block, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := target.(*ArrayInstance)

		newArr := make([]object.EmeraldValue, len(arr.Value))

		for i, val := range arr.Value {
			newArr[i] = yield(block, val)
		}

		return NewArray(newArr)
	}
}

func arrayEach() object.BuiltInMethod {
	return func(target object.EmeraldValue, block *object.Block, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := target.(*ArrayInstance)

		for _, val := range arr.Value {
			yield(block, val)
		}

		return arr
	}
}
