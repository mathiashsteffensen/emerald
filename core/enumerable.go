package core

import (
	"emerald/object"
)

func (rt *Runtime) InitEnumerable() {
	rt.Enumerable = rt.DefineModule("Enumerable")

	rt.DefineMethod(rt.Enumerable, "first", rt.enumerableFirst())
	rt.DefineMethod(rt.Enumerable, "find", rt.enumerableFind())
	rt.DefineMethod(rt.Enumerable, "detect", rt.enumerableFind())
	rt.DefineMethod(rt.Enumerable, "find_index", rt.enumerableFindIndex())
	rt.DefineMethod(rt.Enumerable, "map", rt.enumerableMap())
	rt.DefineMethod(rt.Enumerable, "reduce", rt.enumerableReduce())
	rt.DefineMethod(rt.Enumerable, "inject", rt.enumerableReduce())
	rt.DefineMethod(rt.Enumerable, "sum", rt.enumerableSum())
}

func (rt *Runtime) enumerableFirst() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 0, 1); err != nil {
			return object.NewHeapObject(err)
		}

		var numElements = int64(1)

		if len(args) == 1 {
			val, err := rt.EnforceIntegerArg(args[0])
			if err != nil {
				return object.NewHeapObject(err)
			}

			numElements = val
		}

		var values []object.EmeraldValue
		wrappedBlock := object.NewHeapObject(&object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				// TODO: This doesn't stop iterating after the first element has been found, should probably implement a break keyword or something
				if int64(len(values)) != numElements {
					values = append(values, args[0])
				}
				return rt.NULL
			},
		})

		rt.Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		if numElements == 1 {
			if len(values) == 0 {
				return rt.NULL
			}
			return values[0]
		} else {
			return rt.NewArray(values)
		}
	}
}

func (rt *Runtime) enumerableMap() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		arr := make([]object.EmeraldValue, 0)
		block := ctx.Block

		wrappedBlock := object.NewHeapObject(&object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				arr = append(
					arr,
					rt.EvalBlock(block.Heap.(*object.ClosedBlock), kwargs, args...),
				)
				return rt.NULL
			},
		})

		rt.Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		return rt.NewArray(arr)
	}
}

func (rt *Runtime) enumerableFind() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var firstTruthyElement object.EmeraldValue
		found := false
		block := ctx.Block

		wrappedBlock := object.NewHeapObject(&object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				if found {
					return rt.NULL
				}

				if rt.IsTruthy(rt.EvalBlock(block.Heap.(*object.ClosedBlock), kwargs, args...)) {
					if len(args) < 2 {
						firstTruthyElement = args[0]
					} else {
						firstTruthyElement = rt.NewArray(args)
					}
					found = true
				}
				return rt.NULL
			},
		})

		rt.Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		if !found {
			return rt.NULL
		} else {
			return firstTruthyElement
		}
	}
}

func (rt *Runtime) enumerableFindIndex() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		index, found := 0, false
		block := ctx.Block

		wrappedBlock := object.NewHeapObject(&object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				if found {
					return rt.NULL
				}

				if rt.IsTruthy(rt.EvalBlock(block.Heap.(*object.ClosedBlock), kwargs, args...)) {
					found = true
					return rt.NULL
				}

				index++

				return rt.NULL
			},
		})

		rt.Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		if found {
			return rt.NewInteger(int64(index))
		} else {
			return rt.NewInteger(-1)
		}
	}
}

func (rt *Runtime) enumerableSum() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var accumulated object.EmeraldValue

		blockGiven := ctx.BlockGiven()
		block := ctx.Block

		if len(args) != 0 {
			accumulated = args[0]
		} else {
			accumulated = rt.NewInteger(0)
		}

		wrappedBlock := object.NewHeapObject(&object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				if blockGiven {
					accumulated = rt.Send(accumulated, "+", rt.NULL, map[string]object.EmeraldValue{}, rt.EvalBlock(block.Heap.(*object.ClosedBlock), kwargs, args...))
				} else {
					accumulated = rt.Send(accumulated, "+", rt.NULL, kwargs, args...)
				}

				return rt.NULL
			},
		})

		rt.Send(ctx.Self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		return accumulated
	}
}

// https://apidock.com/ruby/rt.Enumerable/reduce
// TODO: Support symbol argument
func (rt *Runtime) enumerableReduce() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var accumulated object.EmeraldValue

		self := ctx.Self
		block := ctx.Block

		argGiven := len(args) != 0
		if argGiven {
			accumulated = args[0]
		} else {
			accumulated = rt.Send(self, "first", rt.NULL, map[string]object.EmeraldValue{})
		}

		passedFirst := false

		wrappedBlock := object.NewHeapObject(&object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				if argGiven || passedFirst {
					args = append([]object.EmeraldValue{accumulated}, args...)
					accumulated = rt.EvalBlock(block.Heap.(*object.ClosedBlock), kwargs, args...)
				} else {
					passedFirst = true
				}

				return rt.NULL
			},
		})

		rt.Send(self, "each", wrappedBlock, map[string]object.EmeraldValue{})

		return accumulated
	}
}
