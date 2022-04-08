package object

var String *Class

type StringInstance struct {
	*Instance
	Value string
}

func (s *StringInstance) Inspect() string { return s.Value }

func NewString(val string) EmeraldValue {
	return &StringInstance{String.New(), val}
}

func init() {
	String = NewClass(
		"String",
		Object,
		BuiltInMethodSet{
			"to_s": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return target
			},
			"+": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return NewString(target.(*StringInstance).Value + args[0].(*StringInstance).Value)
			},
		},
		BuiltInMethodSet{},
	)
}
