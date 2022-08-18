package core

import "emerald/object"

var Module *object.Class

func InitModule() {
	Module = object.NewClass("Module", Object, Object.Class(), object.BuiltInMethodSet{
		"define_method": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
			ctx.DefinitionTarget.DefineMethod(block, args...)

			return args[0]
		},
	}, object.BuiltInMethodSet{})

	Class.SetSuper(Module)
	Class.Class().(*object.SingletonClass).SetSuper(Module.Class())

	Kernel.Class().(*object.SingletonClass).SetSuper(Module)
}
