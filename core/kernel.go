package core

import (
	"emerald/object"
	"fmt"
)

var Kernel *object.Module

func init() {
	Kernel = object.NewModule(
		"Kernel",
		object.BuiltInMethodSet{
			"class": kernelClass(),
		},
		object.BuiltInMethodSet{
			"puts":    kernelPuts(),
			"include": kernelInclude(),
		},
	)
}

/*** Instance methods ***/

func kernelClass() object.BuiltInMethod {
	return func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return target.ParentClass()
	}
}

/*** Static methods ***/

func kernelPuts() object.BuiltInMethod {
	return func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		strings := []any{}

		for _, arg := range args {
			val, err := arg.SEND(yield, "to_s", arg, nil)
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
	return func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		if len(args) == 0 {
			return NewStandardError("wrong number of arguments (given 0, expected 1+)")
		}

		for _, arg := range args {
			mod, ok := arg.(*object.Module)
			if !ok {
				return NewStandardError(fmt.Sprintf("wrong argument type %s (expected Module)", arg.ParentClass().(*object.Class).Name))
			}

			target.Include(mod)
		}

		return target
	}
}
