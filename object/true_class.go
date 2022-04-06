package object

var TrueClass *Class

var TRUE EmeraldValue

func init() {
	TrueClass = NewClass("TrueClass", Object, BuiltInMethodSet{
		"==": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target != args[0])
		},
	})

	TRUE = TrueClass.New()
}
