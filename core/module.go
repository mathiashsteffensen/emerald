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
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		self.DefineMethod(block, args...)

		return args[0]
	}
}

func moduleAttrReader() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := nameAndInstanceVarFromObject(arg)

			self.BuiltInMethodSet()[name] = func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				value := self.InstanceVariableGet(instanceVarName, self, self)

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
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := nameAndInstanceVarFromObject(arg)

			self.BuiltInMethodSet()[fmt.Sprintf("%s=", name)] = func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				target.InstanceVariableSet(instanceVarName, args[0])

				return args[0]
			}
		}

		return NewArray(args)
	}
}

func moduleAttrAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		_, err := self.SEND(ctx, yield, "attr_reader", self, block, args...)
		if err != nil {
			return NewStandardError(err.Error())
		}

		_, err = self.SEND(ctx, yield, "attr_writer", self, block, args...)
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
