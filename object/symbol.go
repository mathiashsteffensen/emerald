package object

var Symbol *Class

func init() {
	Symbol = NewClass(
		"Symbol",
		Object,
		BuiltInMethodSet{
			"to_s": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
				val := target.Inspect()

				return NewString(val[1:])
			},
		},
	)
}

type SymbolInstance struct {
	*Instance
	Value string
}

func (s *SymbolInstance) Inspect() string { return s.Value }

func NewSymbol(val string) EmeraldValue {
	return &SymbolInstance{Symbol.New(), val}
}
