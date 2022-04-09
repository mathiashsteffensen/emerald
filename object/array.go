package object

var Array *Class

type ArrayInstance struct {
	*Instance
	Value []EmeraldValue
}

var arrayBuiltInMethodSet = BuiltInMethodSet{
	"find": func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue {
		arr := target.(*ArrayInstance)

		for _, val := range arr.Value {
			if yield(block, val) == TRUE {
				return val
			}
		}

		return NULL
	},
	"find_index": func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue {
		arr := target.(*ArrayInstance)

		for i, val := range arr.Value {
			if yield(block, val) == TRUE {
				return NewInteger(int64(i))
			}
		}

		return NULL
	},
	"map": func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue {
		arr := target.(*ArrayInstance)

		newArr := make([]EmeraldValue, len(arr.Value))

		for i, val := range arr.Value {
			newArr[i] = yield(block, val)
		}

		return NewArray(newArr)
	},
	"each": func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue {
		arr := target.(*ArrayInstance)

		for _, val := range arr.Value {
			yield(block, val)
		}

		return arr
	},
	"first": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
		return target.(*ArrayInstance).Value[0]
	},
	"last": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
		targetVal := target.(*ArrayInstance).Value

		return targetVal[len(targetVal)-1]
	},
}

func NewArray(val []EmeraldValue) *ArrayInstance {
	return &ArrayInstance{
		Instance: Array.New(),
		Value:    val,
	}
}
