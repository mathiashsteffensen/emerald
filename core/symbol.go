package core

import "emerald/object"

func (rt *Runtime) InitSymbol() {
	rt.Symbol = rt.DefineClass("Symbol", rt.Object)

	rt.DefineMethod(rt.Symbol, "to_s", rt.symbolToS())
}

type SymbolInstance struct {
	*object.Instance
	Value string
}

func (s *SymbolInstance) Inspect() string { return ":" + s.Value }

func (rt *Runtime) NewSymbol(val string) object.EmeraldValue {
	return rt.GlobalSymbolInternPool.ResolveOrDefine(rt, val)
}

type SymbolInternStore map[string]object.EmeraldValue

func (s SymbolInternStore) Resolve(val string) (object.EmeraldValue, bool) {
	sym, ok := s[val]
	return sym, ok
}

func (s SymbolInternStore) Define(rt *Runtime, val string) object.EmeraldValue {
	sym := &SymbolInstance{Value: val, Instance: rt.Symbol.New()}

	s[val] = sym

	return sym
}

func (s SymbolInternStore) ResolveOrDefine(rt *Runtime, val string) object.EmeraldValue {
	if sym, ok := s.Resolve(val); ok {
		return sym
	} else {
		return s.Define(rt, val)
	}
}

func (rt *Runtime) symbolToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		val := ctx.Self.Inspect()

		return rt.NewString(val[1:])
	}
}
