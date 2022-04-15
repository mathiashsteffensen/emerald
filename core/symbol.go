package core

import "emerald/object"

var Symbol *object.Class

func init() {
	Symbol = object.NewClass(
		"Symbol",
		Object,
		object.BuiltInMethodSet{
			"to_s": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				val := target.Inspect()

				return NewString(val[1:])
			},
		},
		object.BuiltInMethodSet{},
	)
}

type SymbolInstance struct {
	*object.Instance
	Value string
}

func (s *SymbolInstance) Inspect() string { return ":" + s.Value }

func NewSymbol(val string) object.EmeraldValue {
	return &SymbolInstance{Symbol.New(), val}
}
