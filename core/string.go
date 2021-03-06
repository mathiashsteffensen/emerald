package core

import (
	"emerald/object"
	"fmt"
	"strings"
)

var String *object.Class

type StringInstance struct {
	*object.Instance
	Value string
}

func (s *StringInstance) Inspect() string { return s.Value }

func NewString(val string) object.EmeraldValue {
	return &StringInstance{String.New(), val}
}

func InitString() {
	String = object.NewClass(
		"String",
		Object,
		Object.Class(),
		object.BuiltInMethodSet{
			"to_s": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return target
			},
			"to_sym": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewSymbol(target.Inspect())
			},
			"+": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				targetString := target.(*StringInstance)

				argString, ok := args[0].(*StringInstance)
				if !ok {
					var typ string

					if args[0].Type() == object.CLASS_VALUE {
						typ = args[0].(*object.Class).Name
					} else {
						typ = args[0].Class().Super().(*object.Class).Name
					}

					return NewStandardError(fmt.Sprintf("no implicit conversion of %s to String", typ))
				}

				return NewString(targetString.Value + argString.Value)
			},
			"upcase": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewString(strings.ToUpper(target.(*StringInstance).Value))
			},
		},
		object.BuiltInMethodSet{},
	)
}
