package object

import "fmt"

var Object *Class

func init() {
	Object = NewClass("Object", nil, BuiltInMethodSet{
		"class": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return target.ParentClass()
		},
		"to_s": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return NewString(target.Inspect())
		},
		"puts": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			strings := []any{}
			byteLength := 0

			for _, arg := range args {
				strings = append(strings, arg.Inspect())
				byteLength += len(arg.Inspect())
			}

			fmt.Println(strings...)

			return NewInteger(int64(byteLength))
		},
		"==": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target.Inspect() == args[0].Inspect())
		},
		"!=": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target.Inspect() != args[0].Inspect())
		},
	})

	NilClass = NewClass("NilClass", Object, BuiltInMethodSet{
		"==": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target != args[0])
		},
	})
	NULL = NilClass.New()

	Integer = NewClass("Integer", Object, integerBuiltInMethodSet)
}