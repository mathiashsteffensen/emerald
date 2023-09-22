package core

import (
	"emerald/object"
	"fmt"
)

var Module *object.Class

func InitModule() {
	Module = DefineClass("Module", Object)

	DefineMethod(Module, "===", moduleCaseEquals())

	DefineMethod(Module, "define_method", moduleDefineMethod(), object.PRIVATE)
	DefineMethod(Module, "attr_reader", moduleAttrReader(), object.PRIVATE)
	DefineMethod(Module, "attr_writer", moduleAttrWriter(), object.PRIVATE)
	DefineMethod(Module, "attr_accessor", moduleAttrAccessor(), object.PRIVATE)
	DefineMethod(Module, "private", modulePrivate(), object.PRIVATE)

	Class.SetSuper(Module)
	Class.Class().(*object.SingletonClass).SetSuper(Module.Class())

	Kernel.Class().(*object.SingletonClass).SetSuper(Module)
}

func moduleDefineMethod() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		name, err := EnforceArgumentType[*SymbolInstance](Symbol, args[0])
		if err != nil {
			return err
		}

		ctx.Self.DefinedMethodSet()[name.Value] = ctx.Block.(*object.ClosedBlock)

		return args[0]
	}
}

func modulePrivate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if len(args) == 0 {
			ctx.SetDefaultMethodVisibility(object.PRIVATE)
		}

		for _, arg := range args {
		Process:
			switch argTyped := arg.(type) {
			case *StringInstance:
				arg = NewSymbol(argTyped.Value)
				goto Process
			case *SymbolInstance:
				method, _, _, err := ctx.Self.ExtractMethod(argTyped.Value, ctx.Self, ctx.Self)
				if err != nil {
					return Raise(NewStandardError(fmt.Sprintf("undefined method `%s' for class `%s'", argTyped.Value, ctx.Self.Inspect())))
				}

				switch method := method.(type) {
				case *object.ClosedBlock:
					method.Visibility = object.PRIVATE
				case *object.WrappedBuiltInMethod:
					method.Visibility = object.PRIVATE
				}

				return nil
			default:
				return Raise(NewTypeError(fmt.Sprintf("%s is not a symbol nor a string", arg.Inspect())))
			}
		}

		return ctx.Self
	}
}

func moduleAttrReader() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := nameAndInstanceVarFromObject(arg)

			method := func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				value := ctx.Self.InstanceVariableGet(instanceVarName, ctx.Self, ctx.Self)

				if value == nil {
					return NULL
				} else {
					return value
				}
			}

			ctx.Self.BuiltInMethodSet()[name] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}}
		}

		return NewArray(args)
	}
}

func moduleAttrWriter() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := nameAndInstanceVarFromObject(arg)

			method := func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				ctx.Self.InstanceVariableSet(instanceVarName, args[0])

				return args[0]
			}

			ctx.Self.BuiltInMethodSet()[fmt.Sprintf("%s=", name)] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}}
		}

		return NewArray(args)
	}
}

func moduleAttrAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		Send(ctx.Self, "attr_reader", NULL, kwargs, args...)
		Send(ctx.Self, "attr_writer", NULL, kwargs, args...)

		return NULL
	}
}

func moduleCaseEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return Send(args[0], "is_a?", NULL, map[string]object.EmeraldValue{}, ctx.Self)
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
