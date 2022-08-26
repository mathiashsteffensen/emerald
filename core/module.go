package core

import (
	"emerald/object"
	"fmt"
)

var Module *object.Class

func InitModule() {
	Module = object.NewClass(
		"Module",
		Object,
		Object.Class(),
		object.BuiltInMethodSet{
			"define_method": moduleDefineMethod(),
			"attr_reader":   moduleAttrReader(),
			"attr_writer":   moduleAttrWriter(),
			"attr_accessor": moduleAttrAccessor(),
		},
		object.BuiltInMethodSet{},
	)

	Class.SetSuper(Module)
	Class.Class().(*object.SingletonClass).SetSuper(Module.Class())

	Kernel.Class().(*object.SingletonClass).SetSuper(Module)
}

func moduleDefineMethod() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		target.DefineMethod(block, args...)

		return args[0]
	}
}

func moduleAttrReader() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := nameAndInstanceVarFromObject(arg)

			target.BuiltInMethodSet()[name] = func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				value := target.InstanceVariableGet(instanceVarName, target, target)

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
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := nameAndInstanceVarFromObject(arg)

			target.BuiltInMethodSet()[fmt.Sprintf("%s=", name)] = func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				target.InstanceVariableSet(instanceVarName, args[0])

				return args[0]
			}
		}

		return NewArray(args)
	}
}

func moduleAttrAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		_, err := target.SEND(ctx, yield, "attr_reader", target, block, args...)
		if err != nil {
			return NewStandardError(err.Error())
		}

		_, err = target.SEND(ctx, yield, "attr_writer", target, block, args...)
		if err != nil {
			return NewStandardError(err.Error())
		}

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
