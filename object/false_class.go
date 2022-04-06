package object

var FalseClass *Class

var FALSE EmeraldValue

func init() {
	FalseClass = NewClass("FalseClass", Object, BuiltInMethodSet{
		"==": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target != args[0])
		},
	})

	FALSE = FalseClass.New()
}
