package core

import "emerald/object"

var Range *object.Class

type RangeInstance struct {
	*object.Instance
	ExcludeEnd bool
	Begin      object.EmeraldValue
	End        object.EmeraldValue
}

func NewRange(begin object.EmeraldValue, end object.EmeraldValue, excludeEnd bool) *RangeInstance {
	return &RangeInstance{
		Instance:   Range.New(),
		ExcludeEnd: excludeEnd,
		Begin:      begin,
		End:        end,
	}
}

func InitRange() {
	Range = DefineClass(Object, "Range", Object)

	Range.Include(Enumerable)

	DefineSingletonMethod(Range, "new", rangeNew())

	DefineMethod(Range, "each", rangeEach())
}

func rangeNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		begin := args[0]
		end := args[1]
		var excludeEnd bool

		if len(args) < 3 {
			excludeEnd = false
		} else {
			excludeEnd = args[2] == TRUE
		}

		return NewRange(begin, end, excludeEnd)
	}
}

func rangeEach() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		r := ctx.Self.(*RangeInstance)
		begin := r.Begin.(*IntegerInstance).Value
		end := r.End.(*IntegerInstance).Value

		for i := begin; i < end; i++ {
			ctx.Yield(NewInteger(i))
		}

		if !r.ExcludeEnd {
			ctx.Yield(r.End)
		}

		return ctx.Self
	}
}
