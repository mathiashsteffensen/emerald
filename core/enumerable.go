package core

import "emerald/object"

// CRuby docs for Enumerable module https://ruby-doc.org/core-3.1.2/Enumerable.html

var Enumerable *object.Module

func InitEnumerable() {
	Enumerable = object.NewModule(
		"Enumerable",
		object.BuiltInMethodSet{
			"map":        enumerableMap(),
			"find":       enumerableFind(),
			"find_index": enumerableFindIndex(),
			"sum":        enumerableSum(),
		},
		object.BuiltInMethodSet{},
	)
}

func enumerableMap() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		arr := NewArray(make([]object.EmeraldValue, 0))

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, target object.EmeraldValue, _block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				arr.Value = append(
					arr.Value,
					yield(block, args...),
				)
				return NULL
			},
		}

		_, err := target.SEND(ctx, yield, "each", target, wrappedBlock)
		if err != nil {
			return NewStandardError(err.Error())
		}

		return arr
	}
}

func enumerableFind() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		var firstTruthyElement object.EmeraldValue

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, target object.EmeraldValue, _block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				if firstTruthyElement != nil {
					return NULL
				}

				if IsTruthy(yield(block, args...)) {
					if len(args) < 2 {
						firstTruthyElement = args[0]
					} else {
						firstTruthyElement = NewArray(args)
					}
				}
				return NULL
			},
		}

		_, err := target.SEND(ctx, yield, "each", target, wrappedBlock)
		if err != nil {
			return NewStandardError(err.Error())
		}

		if firstTruthyElement == nil {
			return NULL
		} else {
			return firstTruthyElement
		}
	}
}

func enumerableFindIndex() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		index, found := 0, false

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, target object.EmeraldValue, _block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				if found {
					return NULL
				}

				if IsTruthy(yield(block, args...)) {
					found = true
					return NULL
				}

				index++

				return NULL
			},
		}

		_, err := target.SEND(ctx, yield, "each", target, wrappedBlock)
		if err != nil {
			return NewStandardError(err.Error())
		}

		if found {
			return NewInteger(int64(index))
		} else {
			return NewInteger(-1)
		}
	}
}

func enumerableSum() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		var err error
		blockGiven := !IsNull(block)

		var accumulated object.EmeraldValue

		if len(args) != 0 {
			accumulated = args[0]
		} else {
			accumulated = NewInteger(0)
		}

		wrappedBlock := &object.WrappedBuiltInMethod{
			Method: func(ctx *object.Context, target object.EmeraldValue, _block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				if blockGiven {
					accumulated, err = accumulated.SEND(ctx, nil, "+", accumulated, nil, yield(block, args...))
					if err != nil {
						return NewStandardError(err.Error())
					}
				} else {
					accumulated, err = accumulated.SEND(ctx, nil, "+", accumulated, nil, args...)
					if err != nil {
						return NewStandardError(err.Error())
					}
				}

				return NULL
			},
		}

		_, err = target.SEND(ctx, yield, "each", target, wrappedBlock)
		if err != nil {
			return NewStandardError(err.Error())
		}

		return accumulated
	}
}
