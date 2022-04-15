package core

import "emerald/object"

var BuiltIns = map[string]object.EmeraldValue{
	"Object":     Object,
	"NilClass":   NilClass,
	"FalseClass": FalseClass,
	"TrueClass":  TrueClass,
}

func NativeBoolToBooleanObject(input bool) object.EmeraldValue {
	if input {
		return TRUE
	}
	return FALSE
}
