package object

import "fmt"

var Object *Class

func init() {
	Object = NewClass(
		"Object",
		nil,
		BuiltInMethodSet{
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
			"puts": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
				strings := []any{}
				byteLength := 0

				for _, arg := range args {
					strings = append(strings, arg.Inspect())
					byteLength += len(arg.Inspect())
				}

				fmt.Println(strings...)

				return NewInteger(int64(byteLength))
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
}
