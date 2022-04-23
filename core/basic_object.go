package core

import (
	"emerald/object"
)

var BasicObject *object.Class

func InitBasicObject() {
	BasicObject = object.NewClass("BasicObject", nil, object.BuiltInMethodSet{}, object.BuiltInMethodSet{
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
			return target.(*object.StaticClass).Class.(*object.Class).New()
		},
		"define_method": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			ctx.DefinitionTarget.DefineMethod(block, args...)

			return args[0]
		},
	})
}
