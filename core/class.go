package core

import "emerald/object"

var Class *object.Class

func InitClass() {
	Class = object.NewClass(
		"Class",
		nil,
		nil,
		object.BuiltInMethodSet{
			"ancestors": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewArray(ctx.ExecutionTarget.Ancestors())
			},
			"methods": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				methods := []object.EmeraldValue{}

				for _, method := range target.Methods(target) {
					methods = append(methods, NewSymbol(method))
				}

				return NewArray(methods)
			},
			"new": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return ctx.ExecutionTarget.(*object.Class).New()
			},
		},
		object.BuiltInMethodSet{
			"new": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return object.NewClass("", Object, Object.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
			},
		},
	)
}
