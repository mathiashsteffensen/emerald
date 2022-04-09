package object

import (
	"fmt"
	"strings"
)

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
				targetString := target.(*StringInstance)

				argString, ok := args[0].(*StringInstance)
				if !ok {
					var typ string

					if args[0].Type() == CLASS_VALUE {
						typ = args[0].(*Class).Name
					} else {
						typ = args[0].ParentClass().(*Class).Name
					}

					return NewStandardError(fmt.Sprintf("no implicit conversion of %s to String", typ))
				}

				return NewString(targetString.Value + argString.Value)
			},
			"upcase": func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return NewString(strings.ToUpper(target.(*StringInstance).Value))
			},
		},
		BuiltInMethodSet{},
	)
}
