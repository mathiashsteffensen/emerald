package core

import (
	"emerald/object"
	"fmt"
)

var Object *object.Class

var MainObject *object.Instance

func init() {
	Object = object.NewClass(
		"Object",
		nil,
		object.BuiltInMethodSet{
			// Uncomment when 'Object' has been moved to core
			"methods": func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				methods := []object.EmeraldValue{}

				for key := range target.(*object.Instance).DefinedMethodSet() {
					methods = append(methods, NewString(key))
				}

				return NewArray(methods)
			},
			"class": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return target.ParentClass()
			},
			"to_s": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewString(target.Inspect())
			},
			"==": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NativeBoolToBooleanObject(target == args[0])
			},
			"!=": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NativeBoolToBooleanObject(target.Inspect() != args[0].Inspect())
			},
		},
		object.BuiltInMethodSet{
			"new": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return target.(*object.Class).New()
			},
			"define_method": func(target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				target.DefineMethod(false, block, args...)

				return args[0]
			},
			"puts": func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
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
			},
		},
	)

	MainObject = Object.New()
	MainObject.BuiltInSingletonMethods["to_s"] = func(target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString("main:Object")
	}
}
