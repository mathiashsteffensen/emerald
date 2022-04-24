package core

import (
	"emerald/object"
	"fmt"
)

var Kernel *object.Module

func InitKernel() {
	Kernel = object.NewModule(
		"Kernel",
		object.BuiltInMethodSet{
			"class":   kernelClass(),
			"puts":    kernelPuts(),
			"include": kernelInclude(),
		},
		object.BuiltInMethodSet{},
		Module,
	)
}

/*** Instance methods ***/

func kernelClass() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		class := target.Class()
		typ := class.Type()

		for typ != object.CLASS_VALUE {
			class = class.Super()
			typ = class.Type()
		}

		return class
	}
}

/*** Static methods ***/

func kernelPuts() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		strings := []any{}

		for _, arg := range args {
			val, err := arg.SEND(ctx, yield, "to_s", arg, nil)
			if err != nil {
				return NewStandardError(err.Error())
			}

			strings = append(strings, val.Inspect())
		}

		fmt.Println(strings...)

		return NULL
	}
}

func kernelInclude() object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		if len(args) == 0 {
			return NewStandardError("wrong number of arguments (given 0, expected 1+)")
		}

		for _, arg := range args {
			if arg == nil {
				continue
			}

			mod, ok := arg.(*object.Module)
			if !ok {
				return NewStandardError(fmt.Sprintf("wrong argument type %s (expected Module)", arg.Inspect()))
			}

			ctx.DefinitionTarget.Include(mod)
		}

		return ctx.DefinitionTarget
	}
}
