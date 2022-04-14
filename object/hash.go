package object

var Hash = NewClass("Hash", Object, BuiltInMethodSet{}, BuiltInMethodSet{})

type HashInstance struct {
	*Instance
	Value map[EmeraldValue]EmeraldValue
}

func NewHash(val map[EmeraldValue]EmeraldValue) *HashInstance {
	return &HashInstance{
		Instance: Hash.New(),
		Value:    val,
	}
}
