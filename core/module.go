package core

import (
	"emerald/object"
	"fmt"
)

func (rt *Runtime) InitModule() {
	rt.Module = rt.DefineClass("Module", rt.Object)

	rt.DefineMethod(rt.Module, "===", rt.moduleCaseEquals())

	rt.DefineMethod(rt.Module, "name", rt.className())
	rt.DefineMethod(rt.Module, "define_method", rt.moduleDefineMethod(), object.PRIVATE)
	rt.DefineMethod(rt.Module, "attr_reader", rt.moduleAttrReader(), object.PRIVATE)
	rt.DefineMethod(rt.Module, "attr_writer", rt.moduleAttrWriter(), object.PRIVATE)
	rt.DefineMethod(rt.Module, "attr_accessor", rt.moduleAttrAccessor(), object.PRIVATE)
	rt.DefineMethod(rt.Module, "private", rt.modulePrivate(), object.PRIVATE)

	rt.Class.SetSuper(rt.Module)
	rt.Class.Class().(*object.SingletonClass).SetSuper(rt.Module.Class())

	rt.Kernel.Class().(*object.SingletonClass).SetSuper(rt.Module)
}

func (rt *Runtime) moduleDefineMethod() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		name, err := EnforceArgumentType[*SymbolInstance](rt, rt.Symbol, args[0])
		if err != nil {
			return err
		}

		ctx.Self.DefinedMethodSet()[name.Value] = ctx.Block.(*object.ClosedBlock)

		return args[0]
	}
}

func (rt *Runtime) modulePrivate() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if len(args) == 0 {
			ctx.SetDefaultMethodVisibility(object.PRIVATE)
		}

		for _, arg := range args {
		Process:
			switch argTyped := arg.(type) {
			case *StringInstance:
				arg = rt.NewSymbol(argTyped.Value)
				goto Process
			case *SymbolInstance:
				method, _, _, err := ctx.Self.ExtractMethod(argTyped.Value, ctx.Self, ctx.Self)
				if err != nil {
					return rt.Raise(rt.NewStandardError(fmt.Sprintf("undefined method `%s' for class `%s'", argTyped.Value, ctx.Self.Inspect())))
				}

				switch method := method.(type) {
				case *object.ClosedBlock:
					method.Visibility = object.PRIVATE
				case *object.WrappedBuiltInMethod:
					method.Visibility = object.PRIVATE
				}

				return nil
			default:
				return rt.Raise(rt.NewTypeError(fmt.Sprintf("%s is not a symbol nor a string", arg.Inspect())))
			}
		}

		return ctx.Self
	}
}

func (rt *Runtime) moduleAttrReader() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := rt.nameAndInstanceVarFromObject(arg)

			method := func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				value := ctx.Self.InstanceVariableGet(instanceVarName, ctx.Self, ctx.Self)

				if value == nil {
					return rt.NULL
				} else {
					return value
				}
			}

			ctx.Self.BuiltInMethodSet()[name] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}}
		}

		return rt.NewArray(args)
	}
}

func (rt *Runtime) moduleAttrWriter() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			name, instanceVarName := rt.nameAndInstanceVarFromObject(arg)

			method := func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
				ctx.Self.InstanceVariableSet(instanceVarName, args[0])

				return args[0]
			}

			ctx.Self.BuiltInMethodSet()[fmt.Sprintf("%s=", name)] = &object.WrappedBuiltInMethod{Method: method, BaseEmeraldValue: &object.BaseEmeraldValue{}}
		}

		return rt.NewArray(args)
	}
}

func (rt *Runtime) moduleAttrAccessor() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		rt.Send(ctx.Self, "attr_reader", rt.NULL, kwargs, args...)
		rt.Send(ctx.Self, "attr_writer", rt.NULL, kwargs, args...)

		return rt.NULL
	}
}

func (rt *Runtime) moduleCaseEquals() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.Send(args[0], "is_a?", rt.NULL, map[string]object.EmeraldValue{}, ctx.Self)
	}
}

func (rt *Runtime) nameAndInstanceVarFromObject(obj object.EmeraldValue) (string, string) {
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
