package object

var Array *Class

type ArrayInstance struct {
	*Instance
	Value []EmeraldValue
}

func init() {
	Array = NewClass("Array", Object, BuiltInMethodSet{
		"map": func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue {
			arr := target.(*ArrayInstance)

			newArr := make([]EmeraldValue, len(arr.Value))

			for i, val := range arr.Value {
				newArr[i] = yield(target, block, val)
			}

			return NewArray(newArr)
		},
		"first": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
			return target.(*ArrayInstance).Value[0]
		},
	}, BuiltInMethodSet{})
}

func NewArray(val []EmeraldValue) *ArrayInstance {
	return &ArrayInstance{
		Instance: Array.New(),
		Value:    val,
	}
}
