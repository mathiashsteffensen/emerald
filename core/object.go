package core

import (
	"emerald/object"
	"fmt"
)

var Object *object.Class

func init() {
	Object = object.NewClass(
		"Object",
		nil,
		object.BuiltInMethodSet{
			// Uncomment when 'Object' has been moved to core
			/*"methods": func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				methods := []EmeraldValue{}

				for key := range target.(*Instance).DefinedMethodSet() {
					methods = append(methods, NewString(key))
				}

				return NewArray(methods)
			},*/
			"class": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return target.ParentClass()
			},
			"to_s": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewString(target.Inspect())
			},
			"==": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NativeBoolToBooleanObject(target.Inspect() == args[0].Inspect())
			},
			"!=": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NativeBoolToBooleanObject(target.Inspect() != args[0].Inspect())
			},
		},
		object.BuiltInMethodSet{
			"new": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return target.(*object.Class).New()
			},
			"define_method": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				target.DefineMethod(false, block, args...)

				return args[0]
			},
			"puts": func(target object.EmeraldValue, block *object.Block, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				strings := []any{}

				for _, arg := range args {
					strings = append(strings, arg.Inspect())
				}

				fmt.Println(strings...)

				return NULL
			},
		},
	)
}