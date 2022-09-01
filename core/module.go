package core

import (
	"emerald/object"
	"fmt"
)

var Module *object.Class

func InitModule() {
	Module = DefineClass(Object, "Module", Object)

	DefineMethod(Module, "define_method", moduleDefineMethod())
	DefineMethod(Module, "attr_reader", moduleAttrReader())
	DefineMethod(Module, "attr_writer", moduleAttrWriter())
	DefineMethod(Module, "attr_accessor", moduleAttrAccessor())

	Class.SetSuper(Module)
	Class.Class().(*object.SingletonClass).SetSuper(Module.Class())

	Kernel.Class().(*object.SingletonClass).SetSuper(Module)
}

func moduleDefineMethod() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		ctx.Self.DefineMethod(ctx.Block, args...)

		return args[0]
	}
}

func moduleAttrReader() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := nameAndInstanceVarFromObject(arg)

			ctx.Self.BuiltInMethodSet()[name] = func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
				value := ctx.Self.InstanceVariableGet(instanceVarName, ctx.Self, ctx.Self)

				if value == nil {
					return NULL
				} else {
					return value
				}
			}
		}

		return NewArray(args)
	}
}

func moduleAttrWriter() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := nameAndInstanceVarFromObject(arg)

			ctx.Self.BuiltInMethodSet()[fmt.Sprintf("%s=", name)] = func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
				ctx.Self.InstanceVariableSet(instanceVarName, args[0])

				return args[0]
			}
		}

		return NewArray(args)
	}
}

func moduleAttrAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		Send(ctx.Self, "attr_reader", NULL, args...)
		Send(ctx.Self, "attr_writer", NULL, args...)

		return NULL
	}
}

func nameAndInstanceVarFromObject(obj object.EmeraldValue) (string, string) {
	name := ""

	switch obj := obj.(type) {
	case *StringInstance:
		name = obj.Value
	case *SymbolInstance:
		name = obj.Value
	}

	instanceVarName := fmt.Sprintf("@%s", name)

	return name, instanceVarName
}
