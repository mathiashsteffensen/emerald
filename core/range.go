package core

import "emerald/object"

type RangeInstance struct {
	*object.Instance
	ExcludeEnd bool
	Begin      object.EmeraldValue
	End        object.EmeraldValue
}

func (rt *Runtime) NewRange(begin object.EmeraldValue, end object.EmeraldValue, excludeEnd bool) *RangeInstance {
	return &RangeInstance{
		Instance:   rt.Range.New(),
		ExcludeEnd: excludeEnd,
		Begin:      begin,
		End:        end,
	}
}

func (rt *Runtime) InitRange() {
	rt.Range = rt.DefineClass("Range", rt.Object)

	rt.Range.Include(rt.Enumerable)

	rt.DefineSingletonMethod(rt.Range, "new", rt.rangeNew())

	rt.DefineMethod(rt.Range, "each", rt.rangeEach())
}

func (rt *Runtime) rangeNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		begin := args[0]
		end := args[1]
		var excludeEnd bool

		if len(args) < 3 {
			excludeEnd = false
		} else {
			excludeEnd = args[2] == rt.TRUE
		}

		return rt.NewRange(begin, end, excludeEnd)
	}
}

func (rt *Runtime) rangeEach() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		r := ctx.Self.(*RangeInstance)
		begin := r.Begin.(*IntegerInstance).Value
		end := r.End.(*IntegerInstance).Value

		for i := begin; i < end; i++ {
			ctx.Yield(map[string]object.EmeraldValue{}, rt.NewInteger(i))
		}

		if !r.ExcludeEnd {
			ctx.Yield(map[string]object.EmeraldValue{}, r.End)
		}

		return ctx.Self
	}
}
