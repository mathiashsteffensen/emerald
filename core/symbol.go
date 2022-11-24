package core

import "emerald/object"

var Symbol *object.Class

func InitSymbol() {
	Symbol = DefineClass("Symbol", Object)

	DefineMethod(Symbol, "to_s", symbolToS())
}

type SymbolInstance struct {
	*object.Instance
	Value string
}

func (s *SymbolInstance) Inspect() string { return ":" + s.Value }

func NewSymbol(val string) object.EmeraldValue {
	return GlobalSymbolInternPool.ResolveOrDefine(val)
}

type SymbolInternStore map[string]object.EmeraldValue

var GlobalSymbolInternPool = SymbolInternStore{}

func (s SymbolInternStore) Resolve(val string) (object.EmeraldValue, bool) {
	sym, ok := s[val]
	return sym, ok
}

func (s SymbolInternStore) Define(val string) object.EmeraldValue {
	sym := &SymbolInstance{Value: val, Instance: Symbol.New()}

	s[val] = sym

	return sym
}

func (s SymbolInternStore) ResolveOrDefine(val string) object.EmeraldValue {
	if sym, ok := s.Resolve(val); ok {
		return sym
	} else {
		return s.Define(val)
	}
}

func symbolToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		val := ctx.Self.Inspect()

		return NewString(val[1:])
	}
}
