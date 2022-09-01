package core

import (
	"emerald/object"
)

// CRuby docs for Enumerable module https://ruby-doc.org/core-3.1.2/Enumerable.html

var Enumerable *object.Module

func InitEnumerable() {
	Enumerable = DefineModule(Object, "Enumerable")

	DefineMethod(Enumerable, "map", enumerableMap())
	DefineMethod(Enumerable, "find", enumerableFind())
	DefineMethod(Enumerable, "find_index", enumerableFindIndex())
	DefineMethod(Enumerable, "sum", enumerableSum())
}

func enumerableMap() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		arr := NewArray(make([]object.EmeraldValue, 0))
		block := ctx.Block

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
				arr.Value = append(
					arr.Value,
					object.EvalBlock(block.(*object.ClosedBlock), args...),
				)
				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock)

		return arr
	}
}

func enumerableFind() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		var firstTruthyElement object.EmeraldValue
		block := ctx.Block

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
				if firstTruthyElement != nil {
					return NULL
				}

				if IsTruthy(object.EvalBlock(block.(*object.ClosedBlock), args...)) {
					if len(args) < 2 {
						firstTruthyElement = args[0]
					} else {
						firstTruthyElement = NewArray(args)
					}
				}
				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock)

		if firstTruthyElement == nil {
			return NULL
		} else {
			return firstTruthyElement
		}
	}
}

func enumerableFindIndex() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		index, found := 0, false
		block := ctx.Block

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
				if found {
					return NULL
				}

				if IsTruthy(object.EvalBlock(block.(*object.ClosedBlock), args...)) {
					found = true
					return NULL
				}

				index++

				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock)

		if found {
			return NewInteger(int64(index))
		} else {
			return NewInteger(-1)
		}
	}
}

func enumerableSum() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		var accumulated object.EmeraldValue

		blockGiven := ctx.BlockGiven()
		block := ctx.Block

		if len(args) != 0 {
			accumulated = args[0]
		} else {
			accumulated = NewInteger(0)
		}

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
				if blockGiven {
					accumulated = Send(accumulated, "+", NULL, object.EvalBlock(block.(*object.ClosedBlock), args...))
				} else {
					accumulated = Send(accumulated, "+", NULL, args...)
				}

				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock)

		return accumulated
	}
}
