package core

import (
	"emerald/object"
)

var Enumerable *object.Module

func InitEnumerable() {
	Enumerable = DefineModule("Enumerable")

	DefineMethod(Enumerable, "first", enumerableFirst())
	DefineMethod(Enumerable, "find", enumerableFind())
	DefineMethod(Enumerable, "detect", enumerableFind())
	DefineMethod(Enumerable, "find_index", enumerableFindIndex())
	DefineMethod(Enumerable, "map", enumerableMap())
	DefineMethod(Enumerable, "reduce", enumerableReduce())
	DefineMethod(Enumerable, "inject", enumerableReduce())
	DefineMethod(Enumerable, "sum", enumerableSum())
}

func enumerableFirst() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 0, 1); err != nil {
			return err
		}

		var numElements = int64(1)

		if len(args) == 1 {
			arg, err := EnforceArgumentType[*IntegerInstance](Integer, args[0])
			if err != nil {
				return err
			}

			numElements = arg.Value
		}

		var values []object.EmeraldValue
		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				// TODO: This doesn't stop iterating after the first element has been found, should probably implement a break keyword or something
				if int64(len(values)) != numElements {
					values = append(values, args[0])
				}
				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		if numElements == 1 {
			return values[0]
		} else {
			return NewArray(values)
		}
	}
}

func enumerableMap() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		arr := make([]object.EmeraldValue, 0)
		block := ctx.Block

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				arr = append(
					arr,
					object.EvalBlock(block.(*object.ClosedBlock), kwargs, args...),
				)
				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		return NewArray(arr)
	}
}

func enumerableFind() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var firstTruthyElement object.EmeraldValue
		block := ctx.Block

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				if firstTruthyElement != nil {
					return NULL
				}

				if IsTruthy(object.EvalBlock(block.(*object.ClosedBlock), kwargs, args...)) {
					if len(args) < 2 {
						firstTruthyElement = args[0]
					} else {
						firstTruthyElement = NewArray(args)
					}
				}
				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		if firstTruthyElement == nil {
			return NULL
		} else {
			return firstTruthyElement
		}
	}
}

func enumerableFindIndex() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		index, found := 0, false
		block := ctx.Block

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				if found {
					return NULL
				}

				if IsTruthy(object.EvalBlock(block.(*object.ClosedBlock), kwargs, args...)) {
					found = true
					return NULL
				}

				index++

				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		if found {
			return NewInteger(int64(index))
		} else {
			return NewInteger(-1)
		}
	}
}

func enumerableSum() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var accumulated object.EmeraldValue

		blockGiven := ctx.BlockGiven()
		block := ctx.Block

		if len(args) != 0 {
			accumulated = args[0]
		} else {
			accumulated = NewInteger(0)
		}

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				if blockGiven {
					accumulated = Send(accumulated, "+", NULL, map[string]object.EmeraldValue{}, object.EvalBlock(block.(*object.ClosedBlock), kwargs, args...))
				} else {
					accumulated = Send(accumulated, "+", NULL, kwargs, args...)
				}

				return NULL
			},
		}

		Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		return accumulated
	}
}

// https://apidock.com/ruby/Enumerable/reduce
// TODO: Support symbol argument
func enumerableReduce() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var accumulated object.EmeraldValue

		self := ctx.Self
		block := ctx.Block

		argGiven := len(args) != 0
		if argGiven {
			accumulated = args[0]
		} else {
			accumulated = Send(self, "first", NULL, map[string]object.EmeraldValue{})
		}

		passedFirst := false

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				if argGiven || passedFirst {
					args = append([]object.EmeraldValue{accumulated}, args...)
					accumulated = object.EvalBlock(block.(*object.ClosedBlock), kwargs, args...)
				} else {
					passedFirst = true
				}

				return NULL
			},
		}

		Send(self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		return accumulated
	}
}
