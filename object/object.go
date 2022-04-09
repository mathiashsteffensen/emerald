package object

import "fmt"

var Object *Class

func init() {
	Object = NewClass(
		"Object",
		nil,
		BuiltInMethodSet{
			"methods": func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				methods := []EmeraldValue{}

				for key, _ := range target.(*Instance).DefinedMethodSet() {
					methods = append(methods, NewString(key))
				}

				return NewArray(methods)
			},
			"class": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return target.ParentClass()
			},
			"to_s": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return NewString(target.Inspect())
			},
			"==": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return nativeBoolToBooleanObject(target.Inspect() == args[0].Inspect())
			},
			"!=": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return nativeBoolToBooleanObject(target.Inspect() != args[0].Inspect())
			},
		},
		BuiltInMethodSet{
			"new": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return target.(*Class).New()
			},
			"define_method": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				return target.DefineMethod(false, block, args...)
			},
			"puts": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				strings := []any{}

				for _, arg := range args {
					strings = append(strings, arg.Inspect())
				}

				fmt.Println(strings...)

				return NULL
			},
		},
	)

	NilClass = NewClass("NilClass", Object, BuiltInMethodSet{
		"==": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target != args[0])
		},
	}, BuiltInMethodSet{})
	NULL = NilClass.New()

	Integer = NewClass("Integer", Object, integerBuiltInMethodSet, BuiltInMethodSet{})
	Array = NewClass("Array", Object, arrayBuiltInMethodSet, BuiltInMethodSet{})
}
