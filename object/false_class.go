package object

var FalseClass *Class

var FALSE EmeraldValue

func init() {
	FalseClass = NewClass("FalseClass", Object, BuiltInMethodSet{
		"==": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target == args[0])
		},
		"!=": func(target EmeraldValue, block *Block, _yield YieldFunc, args ...EmeraldValue) EmeraldValue {
			return nativeBoolToBooleanObject(target != args[0])
		},
	}, BuiltInMethodSet{})

	FALSE = FalseClass.New()
}
